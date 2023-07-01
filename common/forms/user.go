package forms

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
