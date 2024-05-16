package models

import (
	"gorm.io/gorm"
)

type Docker struct {
	gorm.Model
	ImageName string `gorm:"type:varchar(255);not null;"`
	Tag       string `gorm:"type:varchar(30);not null"`   // 镜像Tag
	ImageId   string `gorm:"type:varchar(255);not null;"` // 镜像id

	UserId uint `gorm:"not null;index"`           // 所属用户id
	Kind   uint `gorm:"default:2;not null;index"` // 种类 1:public 2:private
	Status uint `gorm:"default:2;not null;index"` // 状态 1:已上传到仓库 2:上传中 3:上传失败

	Description string `gorm:"type:varchar(256)"`
	Size        string //镜像大小
}
