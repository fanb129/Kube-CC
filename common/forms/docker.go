package forms

type PullImage struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Tag      string `form:"tag" json:"tag" binding:"required"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type TargetImage struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Tag         string `form:"tag" json:"tag" binding:"required"`
	Description string `form:"description" json:"description"`
	Uid         uint   `form:"uid" json:"uid" binding:"required"`
	Kind        uint   `form:"kind" json:"kind" binding:"required"`
}

type PullImageForm struct {
	PullImage   PullImage   `json:"pull_image" form:"pull_image" binding:"required"`
	TargetImage TargetImage `json:"target_image" form:"target_image" binding:"required"`
}

type SaveImageForm struct {
	ContainerID string      `json:"container_id" form:"container_id" binding:"required"`
	NodeIp      string      `json:"node_ip" form:"node_ip" binding:"required"`
	TargetImage TargetImage `json:"target_image" form:"target_image" binding:"required"`
}

type UpdateImageForm struct {
	Id          uint   `form:"id" json:"id" binding:"required"`
	Kind        uint   `form:"kind" json:"kind"`
	Description string `form:"description" json:"description"`
}
