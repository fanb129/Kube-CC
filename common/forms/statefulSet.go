package forms

type StatefulSetAddForm struct {
	Name      string            `json:"name" form:"name" binding:"required,min=3,max=16"`
	Namespace string            `json:"namespace" form:"namespace" binding:"required"`
	Replicas  int32             `json:"replicas" form:"replicas" binding:"required,gte=1"`
	Image     string            `json:"image" form:"image" binding:"required"`
	Command   []string          `json:"command" form:"command"` // 执行命令 可选
	Args      []string          `json:"args" form:"args"`       // 命令参数 可选
	Ports     []int32           `json:"ports" form:"ports"`     // 端口，可选
	Env       map[string]string `json:"env" form:"env"`         // 环境变量，可选
	ApplyResources
}
