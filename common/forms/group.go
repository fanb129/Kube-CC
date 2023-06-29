package forms

type TransAdmin struct {
	//  Username string `forms:"username" json:"username" binding:"required,min=1,max=16"`
	GroupID    uint `forms:"groupid" json:"groupid" binding:"required"`
	OldAdminID uint `forms:"oldadminid" json:"oldadminid" binding:"required"`
	NewAdminID uint `forms:"newadminid" json:"newadminid" binding:"required"`
}

type RemoveUser struct {
	//  Username string `forms:"username" json:"username" binding:"required,min=1,max=16"`
	UserID uint `forms:"userid" json:"userid" binding:"required"`
}

type AddUser struct {
	GroupID uint
	UserID  uint `forms:"userid" json:"userid" binding:"required"`
}

type GroupUpdateForm struct {
	Name        string `forms:"name" json:"name" binding:"required,min=1,max=16"`
	Description string `forms:"description" json:"description" binding:"required,min=1,max=100"`
}
