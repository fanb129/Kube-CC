package models

import (
	"gorm.io/gorm"
	"time"
)

type Ttl struct {
	gorm.Model
	Namespace   string    `gorm:"type:varchar(50);not null;unique;uniqueIndex:idx_namespace"`
	ExpiredTime time.Time `gorm:"not null;"` // 过期时间
}
