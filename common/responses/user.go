package responses

type UserInfo struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Role        uint   `json:"role"`
	Avatar      string `json:"avatar"`
	Gid         uint   `json:"gid"`
	Group       string `json:"group"`
	Cpu         string `json:"cpu"`
	Memory      string `json:"memory"`
	Storage     string `json:"storage"`
	PvcStorage  string `json:"pvcstorage"`
	Gpu         string `json:"gpu"`
	ExpiredTime string `json:"expired_time"`
	Email       string `json:"email"`
}

type LoginResponse struct {
	Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type CaptchaResponse struct {
	Response
	CaptchaId string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}

type UserInfoResponse struct {
	Response
	UserInfo UserInfo `json:"user_info"`
}

//type UserListResponse struct {
//	Response
//	Page     int        `json:"page"`
//	Total    int        `json:"total"`
//	UserList []UserInfo `json:"user_list"`
//}

type UserListResponse struct {
	Response
	Length   int        `json:"length"`
	UserList []UserInfo `json:"user_list"`
}
