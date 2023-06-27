package forms

import "time"

type NsAddForm struct {
	Uid         uint       `forms:"u_id" json:"u_id" binding:"gte=0"`
	Name        string     `forms:"name" json:"name" binding:"required,min=3,max=16"`
	ExpiredTime *time.Time `forms:"expired_time" json:"expired_time"`
	Cpu         string     `json:"cpu" forms:"cpu" binding:"required"`
	Memory      string     `json:"memory" forms:"memory" binding:"required"`
	Num         int        `json:"num" forms:"num" binding:"required"`
}
