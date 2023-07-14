package forms

type LinuxAddForm struct {
	Kind      uint   `form:"kind" json:"kind" binding:"required,gt=0"`
	Name      string `json:"name" form:"name" binding:"required,min=3,max=16"`
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	ApplyResources
}

// LinuxUpdateForm
//type BatchLinuxAddForm struct {
//	Uid         []uint     `form:"u_id" json:"u_id" binding:"required"`
//	Kind        uint       `form:"kind" json:"kind" binding:"required,gt=0"`
//	ExpiredTime *time.Time `form:"expired_time" json:"expired_time"`
//	Resources
//}

type LinuxUpdateForm struct {
	Name      string `json:"name" form:"name" binding:"required,min=3,max=16"`
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	ApplyResources
}
