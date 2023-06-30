package responses

import "time"

type ImageInfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageId   string    `json:"image_id"`
	UserId    uint      `json:"user_id"`
	Kind      int       `json:"kind"`
}

type CreateResponse struct {
	Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UpdateTagResponse struct {
}

type PullingResponse struct {
	Response
	ImageInfo ImageInfo `json:"image_pullinginfo"`
}

type RemoveResponse struct {
	Response
	ID     string `json:"id"`
	Status uint   `json:"status"`
}

type SaveResponse struct {
}

type ImageInfoResponse struct {
	Response
	ImageInfo ImageInfo `json:"image_info"`
}
type ImageListResponse struct {
	Response
	Page      int         `json:"page"`
	Total     int         `json:"total"`
	ImageList []ImageInfo `json:"user_list"`
}

var NoSuchImage = Response{StatusCode: -1, StatusMsg: "该镜像不存在"}
