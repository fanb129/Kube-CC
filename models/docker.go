package models

import (
	"gorm.io/gorm"
)

type Docker struct {
	gorm.Model
	ImageName string `gorm:"type:varchar(255);not null;"`
	ImageId   string `gorm:"type:varchar(255);not null;"` // 镜像id
	UserId    uint   `gorm:"not null;index"`              // 所属用户id
	Kind      int    `gorm:"default:2;not null"`          // 种类 1:public 2:private
	Tag       string `gorm:"type:varchar(30);not null"`   // 镜像Tag
	Size      string //镜像大小
}
