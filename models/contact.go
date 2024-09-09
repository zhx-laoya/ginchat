// 人员关系表
package models

import (
	"ginchat/utils"

	"gorm.io/gorm"
)
//每个人对应于其他所有对象的关系 
//实现群聊，应该找出TargetId为i的所有OwnerId，从而得到所有的人。
type Contact struct {
	gorm.Model
	OwnerId  uint   //谁的关系信息
	TargetId uint   //对应的谁/群 ID 可以是人的id或群的id
	Type     int    //是属于好友关系还是群聊关系对应关系 1好友 2群 3xx
	Desc     string //预留字段
}

func (table *Contact) TableName() string {
	return "contact"
}

// 根据他的id找到他的所有好友
func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id =? and type =1", userId).Find(&contacts)
	for _, v := range contacts { //把表中和他有关系的人找出来放入obj数组中
		objIds = append(objIds, uint64(v.TargetId))
	}
	//根据obj中的所有用户id找出所有用户对象
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

// 添加好友 自己的ID ,好友的名字，根据名称添加用户，返回(-1,0)，以及message
func AddFriend(userId uint, targetName string) (int, string) {
	//用户名不为空
	if targetName != "" {
		targetUser := FindUserByName(targetName)
		//用户注册过，即用户存在
		if targetUser.Salt != "" {
			if targetUser.ID == userId {
				return -1, "不能加自己"
			}
			contact0 := Contact{}
			utils.DB.Where("owner_id=? and target_id=? and type=1", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "不能重复添加"
			}
			//事务一旦开始，不论什么异常最终都会Rollback，事务是为了双向添加好友
			tx := utils.DB.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetUser.ID
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetUser.ID
			contact1.TargetId = userId
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()
			return 0, "添加好友成功"
		} else {
			return -1, "没有找到此用户"
		}
	}
	return -1, "好友名称不能为空"
}

// 加群聊
func JoinGroups(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	contact.Type = 2
	community := Community{}
	utils.DB.Where("id = ? ", comId).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群聊"
	}
	utils.DB.Where("owner_id=? and target_id=? and type=2", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}

// 找到这个群中的所有用户的id
func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
