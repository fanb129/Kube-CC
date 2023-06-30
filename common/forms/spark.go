package forms

import "time"

type SparkAddForm struct {
	Uid            uint       `forms:"u_id" json:"u_id" binding:"required,gt=0"`
	MasterReplicas int32      `forms:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32      `forms:"worker_replicas" json:"worker_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime    *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}

// BatchSparkAddForm 支持批量添加
type BatchSparkAddForm struct {
	Uid            []uint     `forms:"u_id" json:"u_id" binding:"required"`
	MasterReplicas int32      `forms:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32      `forms:"worker_replicas" json:"worker_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime    *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}
type SparkUpdateForm struct {
	Name           string     `json:"name" forms:"name" binding:"required"`
	Uid            uint       `forms:"u_id" json:"u_id" binding:"required,gt=0"`
	MasterReplicas int32      `forms:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32      `forms:"worker_replicas" json:"worker_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime    *time.Time `forms:"expired_time" json:"expired_time"`
	Resources
}
