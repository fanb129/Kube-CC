package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null;unique;uniqueIndex:idx_group"` // 组名
	Adminid     uint   `gorm:"not null"`                                               // 组管理员id
	Description string `gorm:"type:varchar(100)"`                                      // 组的简介
}
