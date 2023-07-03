package forms

type PvcAddForm struct {
	Name             string `json:"name" form:"name" binding:"required"`
	Namespace        string `json:"namespace" form:"namespace" binding:"required,min=3,max=16"`
	StorageClassName string `json:"storage_class_name" form:"storage_class_name" binding:"required"`
	StorageSize      string `json:"storage_size" form:"storage_size" binding:"required"`
	AccessModes      string `json:"access_modes" form:"access_modes" binding:"required"`
}

type PvcUpdateForm struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Namespace   string `json:"namespace" form:"namespace" binding:"required,min=3,max=16"`
	StorageSize string `json:"storage_size" form:"storage_size" binding:"required"`
}
