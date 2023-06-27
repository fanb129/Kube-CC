package forms

import "time"

type HadoopAddForm struct {
	Uid                uint       `forms:"u_id" json:"u_id" binding:"required,gt=0"`
	HdfsMasterReplicas int32      `forms:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32      `forms:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32      `forms:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32      `forms:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime        *time.Time `forms:"expired_time" json:"expired_time"`
	Cpu                string     `json:"cpu" forms:"cpu" binding:"required"`
	Memory             string     `json:"memory" forms:"memory" binding:"required"`
}

// BatchHadoopAddForm 批量添加
type BatchHadoopAddForm struct {
	Uid                []uint     `forms:"u_id" json:"u_id" binding:"required"`
	HdfsMasterReplicas int32      `forms:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32      `forms:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32      `forms:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32      `forms:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime        *time.Time `forms:"expired_time" json:"expired_time"`
	Cpu                string     `json:"cpu" forms:"cpu" binding:"required"`
	Memory             string     `json:"memory" forms:"memory" binding:"required"`
}

type HadoopUpdateForm struct {
	Name               string     `json:"name" forms:"name" binding:"required"`
	Uid                uint       `forms:"u_id" json:"u_id" binding:"required,gt=0"`
	HdfsMasterReplicas int32      `forms:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32      `forms:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32      `forms:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32      `forms:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime        *time.Time `forms:"expired_time" json:"expired_time"`
	Cpu                string     `json:"cpu" forms:"cpu" binding:"required"`
	Memory             string     `json:"memory" forms:"memory" binding:"required"`
}
