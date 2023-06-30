package forms

import "time"

type LinuxAddForm struct {
	Uid         uint       `forms:"u_id" json:"u_id" binding:"required,gt=0"`
	Kind        uint       `forms:"kind" json:"kind" binding:"required,gt=0"`
	ExpiredTime *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}
type BatchLinuxAddForm struct {
	Uid         []uint     `forms:"u_id" json:"u_id" binding:"required"`
	Kind        uint       `forms:"kind" json:"kind" binding:"required,gt=0"`
	ExpiredTime *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}
type LinuxUpdateForm struct {
	Name        string     `json:"name" forms:"name" binding:"required"`
	Uid         uint       `forms:"u_id" json:"u_id" binding:"required,gt=0"`
	ExpiredTime *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}
