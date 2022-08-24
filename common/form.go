package common

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

type ResetPassForm struct {
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

type EditForm struct {
	Role uint `form:"role" json:"role" binding:"required,oneof=1 2 3"`
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
