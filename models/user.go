package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;unique;uniqueIndex:idx_username"` // 账号
	Password string `gorm:"type:varchar(100);not null;"`                               // 加密后的密码
	Nickname string `gorm:"type:varchar(20);not null;index"`                           // 昵称
	Role     uint   `gorm:"not null;default:1"`                                        // 权限状态 1普通用户 2管理员 3超级管理员
	Avatar   string `gorm:"type:varchar(100)"`                                         // 头像地址
}

//func (u *User) Create
