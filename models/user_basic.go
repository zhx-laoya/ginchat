package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Avatar        string //头像
	Identity      string
	ClentIP       string    //设备
	ClientPort    string    //客户端口
	Salt          string    //随机数加密密码
	LoginTime     time.Time //登录时间
	HeartbeatTime time.Time //心跳时间
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"` //下线时间
	IsLogout      bool      //是否下线
	DeviceInfo    string    //设备信息
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// 返回一个可以操控的整个数据库对象，这个是获取所有用户的信息
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 30)
	//注意加&
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

// 根据名字找到用户
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// 根据电话找到用户
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

// 根据email找到用户
func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

// 根据用户名和密码找到用户 每次登录时使用并更新token，用于登录
func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)
	//token加密 ,根据当前时间映射成一个md5字符串
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id= ?", user.ID).Update("identity", temp)
	return user
}

// 传入一个数据让数据库赋值给他
func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

// 传入一个数据 recvProc让数据库删除他
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

// 传入一个数据让数据库修改他
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name,
		PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}

// 根据id查找某个用户
func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id =?", id).First(&user)
	return user
}
