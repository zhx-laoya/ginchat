// 群聊关系表
package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icno    string //图片
	Type    int    //群聊等级
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
