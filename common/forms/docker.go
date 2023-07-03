package forms

type ImageUpdateForm struct {
	ImageId string `form:"imageid" json:"imageid" binding:"required,min=1,max=64"`
}

type SaveForm struct {
	imglist []string `form:"imagelist" json:"imagelist" binding:"required"`
}

type PullSpecifiedForm struct {
	ImageId string `form:"imageid" json:"imageid" binding:"required,min=1,max=64"`
}

type PullFromRepository struct {
	ImageId string `form:"imageid" json:"imageid" binding:"required,min=1,max=64"`
	Account string `form:"account" json:"account" binding:"required,min=1,max=16"`
	Passwd  string `form:"passwd" json:"passwd" binding:"required,min=8,max=64"`
}
