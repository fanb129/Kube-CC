package common

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=1,max=16"`
}

type ResetPassForm struct {
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

type EditForm struct {
	Role uint `form:"role" json:"role" binding:"required,oneof=1 2 3"`
}
type NsAddForm struct {
	Uid  uint   `form:"u_id" json:"u_id" binding:"gte=0"`
	Name string `form:"name" json:"name" binding:"required"`
}
type SparkAddForm struct {
	Uid            uint  `form:"u_id" json:"u_id" binding:"required,gt=0"`
	MasterReplicas int32 `form:"master_replicas" json:"master_replicas" binding:"required,gte=1,lte=3"`
	WorkerReplicas int32 `form:"worker_replicas" json:"worker_replicas" binding:"required,gte=2,lte=10"`
}

type LinuxAddForm struct {
	Uid  uint `form:"u_id" json:"u_id" binding:"required,gt=0"`
	Kind uint `form:"kind" json:"kind" binding:"required,gt=0"`
}

type HadoopAddForm struct {
	Uid                uint  `form:"u_id" json:"u_id" binding:"required,gt=0"`
	HdfsMasterReplicas int32 `form:"hdfs_master_replicas" json:"hdfs_master_replicas" binding:"required,gte=1,lte=3"`
	DatanodeReplicas   int32 `form:"datanode_replicas" json:"datanode_replicas" binding:"required,gte=2,lte=10"`
	YarnMasterReplicas int32 `form:"yarn_master_replicas" json:"yarn_master_replicas" binding:"required,gte=1,lte=3"`
	YarnNodeReplicas   int32 `form:"yarn_node_replicas" json:"yarn_node_replicas" binding:"required,gte=2,lte=10"`
}
