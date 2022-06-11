package models

import "time"

// Model 定义基础模型类，实现复用
type Model struct {
	ID       uint      `gorm:"primaryKey" json:"id"` // id
	CreateAt time.Time `json:"createAt"`             // 创建时间
	UpdateAt time.Time `json:"updateAt"`             // 更新时间
}
