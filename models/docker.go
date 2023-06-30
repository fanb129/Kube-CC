package models

import (
	"gorm.io/gorm"
)

type Docker struct {
	gorm.Model
	ImageId string `gorm:"type:varchar(50);not null;unique;"` // 镜像id
	UserId  uint   `gorm:"not null;index"`                    // 所属用户id
	Kind    int    `gorm:"default:2;not null"`                // 种类 1:public 2:private
}
