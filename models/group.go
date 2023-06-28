package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name    string `gorm:"type:varchar(20);not null;unique;uniqueIndex:idx_username"` // 账号
	AdminId uint   `gorm:"not null"`                                                  // 组管理员id
}
