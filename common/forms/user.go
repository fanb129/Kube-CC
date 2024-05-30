package forms

import "mime/multipart"

type ResetPassForm struct {
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

type EditForm struct {
	Role uint `form:"role" json:"role" binding:"required,oneof=1 2 3"`
}

type UpdateForm struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=1,max=16"` // 昵称
	Avatar   string `form:"avatar" json:"avatar"`
}

type AllocationForm struct {
	Cpu         string `form:"cpu" json:"cpu" binding:"required,min=1,max=16"`
	Memory      string `form:"memory" json:"memory" binding:"required,min=1,max=16"`
	Storage     string `form:"storage" json:"storage" binding:"required,min=1,max=16"`
	Pvcstorage  string `form:"pvcstorage" json:"pvcstorage" binding:"required,min=1,max=16"`
	Gpu         string `form:"gpu" json:"gpu" binding:"required,min=1,max=16"`
	ExpiredTime string `form:"expired_time" json:"expired_time" binding:"required"`
}

type SetEmailForm struct {
	Id    uint   `json:"id" form:"id" binding:"required"`
	Email string `json:"email" form:"email" binding:"required"`
	VCode string `json:"vcode" form:"vcode" binding:"required"`
}

type AddUserForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=1,max=16"`
	Gid      uint   `form:"gid" json:"gid"`
}

type AddUserByFileForm struct {
	File *multipart.FileHeader `form:"file" json:"file" binding:"required"`
	Gid  uint                  `form:"gid" json:"gid"`
}
