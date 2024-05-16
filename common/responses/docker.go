package responses

type ImageInfo struct {
	Id        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	ImageId string `json:"image_id"`
	Name    string `json:"image_name"`
	Tag     string `json:"tag"`
	Size    string `json:"size"`

	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Uid      uint   `json:"u_id"`

	Kind   uint `json:"kind"`
	Status uint `json:"status"`

	Description string `json:"description"`
}

type ImageListResponse struct {
	Response
	Length    int         `json:"length"`
	ImageList []ImageInfo `json:"image_list"`
}

type ImageInfoResponse struct {
	Response
	ImageInfo ImageInfo `json:"image_info"`
}
