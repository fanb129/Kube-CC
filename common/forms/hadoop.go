package forms

import "time"

type HadoopAddForm struct {
	Uid                string     `form:"u_id" json:"u_id" binding:"required"`
	HdfsMasterReplicas int32      `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1"`
	DatanodeReplicas   int32      `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=1"`
	YarnMasterReplicas int32      `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1"`
	YarnNodeReplicas   int32      `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=1"`
	ExpiredTime        *time.Time `form:"expired_time" json:"expired_time"`
	ApplyResources
}

// BatchHadoopAddForm 批量添加
type BatchHadoopAddForm struct {
	Uid                []string   `form:"u_id" json:"u_id" binding:"required"`
	HdfsMasterReplicas int32      `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1"`
	DatanodeReplicas   int32      `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=1"`
	YarnMasterReplicas int32      `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1"`
	YarnNodeReplicas   int32      `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=1"`
	ExpiredTime        *time.Time `form:"expired_time" json:"expired_time"`
	ApplyResources
}

type HadoopUpdateForm struct {
	Name               string     `json:"name" form:"name" binding:"required"`
	HdfsMasterReplicas int32      `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1"`
	DatanodeReplicas   int32      `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=1"`
	YarnMasterReplicas int32      `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1"`
	YarnNodeReplicas   int32      `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=1"`
	ExpiredTime        *time.Time `form:"expired_time" json:"expired_time"`
	ApplyResources
}
