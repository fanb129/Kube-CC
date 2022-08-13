package models

import "gorm.io/gorm"

type Spark struct {
	gorm.Model
	Uid  uint   `gorm:"not null;index"`
	Time string `gorm:"type:varchar(20);not null"`
}
