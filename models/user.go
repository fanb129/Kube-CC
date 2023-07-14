package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;unique;uniqueIndex:idx_username"` // 账号
	Email    string `gorm:"type:varchar(256);not  null;unique;uniqueIndex:idx_email"`  // 邮箱

	Password string `gorm:"type:varchar(100);not null;"`     // 加密后的密码
	Nickname string `gorm:"type:varchar(20);not null;index"` // 昵称
	Avatar   string `gorm:"type:varchar(100)"`               // 头像地址

	Role    uint `gorm:"not null;default:1"` // 权限状态 1普通用户 2组管理员 3超级管理员
	Groupid uint `gorm:"not null;default:0"` // 所属组

	Cpu        string `gorm:"type:varchar(20);not null;default:5"`    //Cpu配额
	Memory     string `gorm:"type:varchar(20);not null;default:10Gi"` //内存配额
	Storage    string `gorm:"type:varchar(20);not null;default:20Gi"` //存储配额
	Pvcstorage string `gorm:"type:varchar(20);not null;default:20Gi"` //持久化存储配额
	Gpu        string `gorm:"type:varchar(20);not null;default:5Gi"`  //Gpu配额

	ExpiredTime time.Time
}

//func (u *User) Create
