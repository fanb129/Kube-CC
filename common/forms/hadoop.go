package forms

import "time"

type HadoopAddForm struct {
	Uid                uint       `form:"u_id" json:"u_id" binding:"required,gt=0"`
	HdfsMasterReplicas int32      `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32      `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32      `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32      `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime        *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}

// BatchHadoopAddForm 批量添加
type BatchHadoopAddForm struct {
	Uid                []uint     `form:"u_id" json:"u_id" binding:"required"`
	HdfsMasterReplicas int32      `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32      `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32      `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32      `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime        *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}

type HadoopUpdateForm struct {
	Name               string     `json:"name" form:"name" binding:"required"`
	Uid                uint       `form:"u_id" json:"u_id" binding:"required,gt=0"`
	HdfsMasterReplicas int32      `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32      `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32      `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32      `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
	ExpiredTime        *time.Time `form:"expired_time" json:"expired_time"`
	Resources
}
