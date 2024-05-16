package forms

type SparkAddForm struct {
	Uid            string `form:"u_id" json:"u_id" binding:"required"`
	Name           string `form:"name" json:"name" binding:"required,min=3,max=16"`
	MasterReplicas int32  `form:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32  `form:"worker_replicas" json:"worker_replicas" binding:"required,gte=1,lte=10"`
	//ExpiredTime    *time.Time `form:"expired_time" json:"expired_time"`
	ApplyResources
}

// BatchSparkAddForm 支持批量添加
type BatchSparkAddForm struct {
	Name           string   `form:"name" json:"name" binding:"required,min=3,max=16"`
	Uid            []string `form:"u_id" json:"u_id" binding:"required"`
	MasterReplicas int32    `form:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32    `form:"worker_replicas" json:"worker_replicas" binding:"required,gte=1,lte=10"`
	//ExpiredTime    *time.Time `form:"expired_time" json:"expired_time"`
	ApplyResources
}
type SparkUpdateForm struct {
	Name           string `json:"name" form:"name" binding:"required"`
	MasterReplicas int32  `form:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32  `form:"worker_replicas" json:"worker_replicas" binding:"required,gte=2,lte=10"`
	//ExpiredTime    *time.Time `form:"expired_time" json:"expired_time"`
	ApplyResources
}
