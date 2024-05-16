package forms

type TransAdmin struct {
	//  Username string `form:"username" json:"username" binding:"required,min=1,max=16"`
	//GroupID    uint `form:"groupid" json:"groupid" binding:"required"`
	OldAdminID uint `form:"oldadminid" json:"oldadminid" binding:"required"`
	NewAdminID uint `form:"newadminid" json:"newadminid" binding:"required"`
}

type RemoveUser struct {
	//  Username string `form:"username" json:"username" binding:"required,min=1,max=16"`
	UserID uint `form:"userid" json:"userid" binding:"required"`
}

type AddUser struct {
	GroupID uint `form:"groupid" json:"groupid" binding:"required"`
	//UserID  uint `form:"userid" json:"userid" binding:"required"`
}

type GroupUpdateForm struct {
	Name        string `form:"name" json:"name" binding:"required,min=1,max=16"`
	Description string `form:"description" json:"description"`
}

// type CreateGroupForm struct {
// 	AdminID     uint   `form:"oldadminid" json:"oldadminid" binding:"required"`
// 	Name        string `form:"name" json:"name" binding:"required,min=1,max=16"`
// 	Description string `form:"description" json:"description" binding:"required,min=1,max=100"`
// }
