package forms

// 修改镜像的tag
type ImageUpdateForm struct {
	OldRepositoryName string `form:"image_old_name"json:"image_old_name" binding:"required,min=1,max=255"`
	OldTag            string `form:"image_old_tag" json:"image_old_tag" binding:"required,min=1,max=64"`
	NewRepositoryName string `form:"image_new_name"json:"image_new_name" binding:"required,min=1,max=255"`
	NewTag            string `form:"image_new_tag" json:"image_new_tag" binding:"required,min=1,max=64"`
}

type SaveForm struct {
	Imglist []string `form:"image_list" json:"image_list" binding:"required"`
}

type PullFromRepositoryPrivate struct {
	RepositoryName string `form:"image_name" json:"image_name" binding:"required,min=1,max=255"`
	Tag            string `form:"tag" json:"tag" binding:"required,min=1,max=64"`
	Uid            uint   `form:"uid" json:"uid"`
	Kind           int    `form:"kind" json:"kind"`
	Username       string `form:"username"json:"username"binding:"required min=1,max=64"`
	Passwd         string `form:"passwd"json:"passwd"binding:"required min=1,max=64"`
}

type PullFromRepositoryPublic struct {
	RepositoryName string `form:"image_name" json:"image_name" binding:"required,min=1,max=255"`
	Tag            string `form:"tag" json:"tag" binding:"required,min=1,max=64"`
	Uid            uint   `form:"uid" json:"uid"`
	Kind           int    `form:"kind" json:"kind"`
}

type ImageCreateForm struct {
	Parent   string `form:"parent" json:"parent" binding:"required,min=1,max=64"`
	Username string `form:"username"json:"username"binding:"required min=1,max=64"`
	Passwd   string `form:"passwd"json:"passwd"binding:"required min=1,max=64"`
	Tag      string `form:"tag" json:"tag" binding:"required,min=1,max=64"`
	Uid      uint   `form:"uid" json:"uid"`
	Kind     int    `form:"kind" json:"kind"`
}
