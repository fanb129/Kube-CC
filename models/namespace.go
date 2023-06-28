package models

import "gorm.io/gorm"

// Namespace 暂时未使用
type Namespace struct {
	gorm.Model
	Uid   uint `gorm:"not null;index"`
	Image string
	Ns    string `gorm:"type:varchar(50);not null;index"`
}
