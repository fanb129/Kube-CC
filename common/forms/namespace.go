package forms

type Resources struct {
	Cpu        string `json:"cpu" form:"cpu" binding:"required"`         //500m = .5 cores
	Memory     string `json:"memory" form:"memory" binding:"required"`   //500Gi 256Mi
	Storage    string `json:"storage" form:"storage" binding:"required"` //临时存储 500Gi = 500GiB
	PvcStorage string `json:"pvc_storage" form:"pvc_storage"`            // pvc持久存储 可选项
	// TODO:GPU
	Gpu string `json:"gpu" form:"gpu"` // GPU 可选项
}

type NsAddForm struct {
	Uid  uint   `form:"u_id" json:"u_id" binding:"gte=1"` // 必须绑定用户
	Name string `form:"name" json:"name" binding:"required,min=3,max=16"`
	//ExpiredTime *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}

type NsUpdateForm struct {
	Name string `form:"name" json:"name" binding:"required"`
	//ExpiredTime *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}
