package forms

import "time"

type Resources struct {
	Cpu        string `json:"cpu" forms:"cpu" binding:"required"`                 //500m = .5 cores
	Memory     string `json:"memory" forms:"memory" binding:"required"`           //500Gi 256Mi
	Storage    string `json:"storage" forms:"storage" binding:"required"`         //临时存储 500Gi = 500GiB
	PvcStorage string `json:"pvc_storage" forms:"pvc_storage" binding:"required"` // pvc存储
	// TODO:GPU
	Gpu string `json:"gpu" forms:"gpu" binding:"required"` // GPU
}

type NsAddForm struct {
	Uid         uint       `forms:"u_id" json:"u_id" binding:"gte=0"`
	Name        string     `forms:"name" json:"name" binding:"required,min=3,max=16"`
	ExpiredTime *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}
