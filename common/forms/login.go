package forms

type LoginForm struct {
	UsernameorEmail string `form:"usernameoremail" json:"usernameoremail" binding:"required,min=3,max=16"`
	Password        string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=1,max=16"`
	Email    string `form:"email" json:"email" binding:"required,min=3,max=16"`
}
