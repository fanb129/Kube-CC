package forms

type JobAddForm struct {
	Name        string   `json:"name" form:"name" binding:"required,min=3,max=16"`
	Namespace   string   `json:"namespace" form:"namespace" binding:"required"`
	Completions int32    `json:"completions" form:"completions" binding:"required,gte=1"` // 指定job需要成功运行Pods的次数。
	Parallelism int32    `json:"parallelism" form:"parallelism" binding:"required,gte=1"` // 指定job在任一时刻应该并发运行Pods的数量
	Image       string   `json:"image" form:"image" binding:"required"`
	Command     []string `json:"command" form:"command"` // 执行命令 可选
	Args        []string `json:"args" form:"args"`       // 命令参数 可选
}
