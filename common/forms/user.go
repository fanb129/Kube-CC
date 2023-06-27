package forms

type ResetPassForm struct {
	Password string `forms:"password" json:"password" binding:"required,min=6,max=16"`
}

type EditForm struct {
	Role uint `forms:"role" json:"role" binding:"required,oneof=1 2 3"`
}

type UpdateForm struct {
	Nickname string `forms:"nickname" json:"nickname" binding:"required,min=1,max=16"` // 昵称
	Avatar   string `forms:"avatar" json:"avatar"`
}
