package forms

type LoginForm struct {
	UsernameorEmail string `forms:"usernameoremail" json:"usernameoremail" binding:"required,min=3,max=16"`
	Password        string `forms:"password" json:"password" binding:"required,min=6,max=16"`
}

type RegisterForm struct {
	Username string `forms:"username" json:"username" binding:"required,min=3,max=16"`
	Password string `forms:"password" json:"password" binding:"required,min=6,max=16"`
	Nickname string `forms:"nickname" json:"nickname" binding:"required,min=1,max=16"`
	Email    string `forms:"email" json:"email" binding:"required,min=3,max=16"`
}
