package forms

type ApplyResources struct {
	Cpu     string `json:"cpu" form:"cpu" binding:"required"`         //500m = .5 cores
	Memory  string `json:"memory" form:"memory" binding:"required"`   //500Gi 256Mi
	Storage string `json:"storage" form:"storage" binding:"required"` //临时存储 500Gi = 500GiB

	PvcStorage       string   `json:"pvc_storage" form:"pvc_storage"` // pvc持久存储 可选项
	StorageClassName string   `json:"storage_class_name" form:"storage_class_name"`
	PvcPath          []string `json:"pvc_path" form:"pvc_path"`
	// TODO:GPU
	Gpu string `json:"gpu" form:"gpu"` // GPU 可选项
}

type DeployAddForm struct {
	Name      string            `json:"name" form:"name" binding:"required,min=3,max=16"`
	Namespace string            `json:"namespace" form:"namespace" binding:"required,min=3,max=16"`
	Replicas  int32             `json:"replicas" form:"replicas" binding:"required,gte=1"`
	Image     string            `json:"image" form:"image" binding:"required"`
	Command   []string          `json:"command" form:"command"` // 执行命令 可选
	Args      []string          `json:"args" form:"args"`       // 命令参数 可选
	Ports     []int32           `json:"ports" form:"ports"`     // 端口，可选
	Env       map[string]string `json:"env" form:"env"`         // 环境变量，可选
	ApplyResources
}
