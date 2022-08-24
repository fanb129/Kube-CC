package models

import "gorm.io/gorm"

type Linux struct {
	gorm.Model
	Uid   uint   `gorm:"not null;index"`
	Image uint   `gorm:"not null"` // 1代表centos，2代表ubuntu
	Time  string `gorm:"type:varchar(20);not null"`
}
