package models

import (
	"time"

	"gorm.io/gorm"
)

type SysLoginLog struct {
	gorm.Model
	Username      string    `json:"username" gorm:"type:varchar(20);not null;unique;uniqueIndex:idx_username"` // 账号
	Status        string    `json:"status" gorm:"type:varchar(10);not null;"`                                  // success or failed
	LoginTime     time.Time `json:"login_time" gorm:"not null;"`
	Ipaddr        string    `json:"ipaddr" gorm:"type:varchar(255);comment:ip地址"`
	LoginLocation string    `json:"loginLocation" gorm:"type:varchar(255);comment:归属地"`
	Browser       string    `json:"browser" gorm:"type:varchar(255);comment:浏览器"`
	Os            string    `json:"os" gorm:"type:varchar(255);comment:系统"`
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
