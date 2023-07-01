package forms

import "time"

type LinuxAddForm struct {
	Uid         uint       `form:"u_id" json:"u_id" binding:"required,gt=0"`
	Kind        uint       `form:"kind" json:"kind" binding:"required,gt=0"`
	ExpiredTime *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}
type BatchLinuxAddForm struct {
	Uid         []uint     `form:"u_id" json:"u_id" binding:"required"`
	Kind        uint       `form:"kind" json:"kind" binding:"required,gt=0"`
	ExpiredTime *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}
type LinuxUpdateForm struct {
	Name        string     `json:"name" form:"name" binding:"required"`
	Uid         uint       `form:"u_id" json:"u_id" binding:"required,gt=0"`
	ExpiredTime *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}
