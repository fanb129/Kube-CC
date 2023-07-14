package models

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/storage"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ActiveRecord interface {
	schema.Tabler
	SetCreateBy(createBy int)
	SetUpdateBy(updateBy int)
	Generate() ActiveRecord
	GetId() interface{}
}
type ControlBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
}

type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type SysLoginLog struct {
	Model
	Username      string    `json:"username" gorm:"size:128;comment:用户名"`
	Status        string    `json:"status" gorm:"size:4;comment:状态"`
	Ipaddr        string    `json:"ipaddr" gorm:"size:255;comment:ip地址"`
	LoginLocation string    `json:"loginLocation" gorm:"size:255;comment:归属地"`
	Browser       string    `json:"browser" gorm:"size:255;comment:浏览器"`
	Os            string    `json:"os" gorm:"size:255;comment:系统"`
	Platform      string    `json:"platform" gorm:"size:255;comment:固件"`
	LoginTime     time.Time `json:"loginTime" gorm:"comment:登录时间"`
	Remark        string    `json:"remark" gorm:"size:255;comment:备注"`
	Msg           string    `json:"msg" gorm:"size:255;comment:信息"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	ControlBy
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}

func (e *SysLoginLog) Generate() ActiveRecord {
	o := *e
	return &o
}

func (e *SysLoginLog) GetId() interface{} {
	return e.Id
}

// SaveLoginLog 从队列中获取登录日志
func SaveLoginLog(message storage.Messager) (err error) {
	//准备db
	db := sdk.Runtime.GetDbByKey(message.GetPrefix())
	if db == nil {
		err = errors.New("db not exist")
		log.Errorf("host[%s]'s %s", message.GetPrefix(), err.Error())
		return err
	}

	var rb []byte

	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var l SysLoginLog
	err = json.Unmarshal(rb, &l)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	err = db.Create(&l).Error
	if err != nil {
		log.Errorf("db create error, %s", err.Error())
		return err
	}
	return nil
}
