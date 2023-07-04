package responses

type UserInfo struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Role      uint   `json:"role"`
	Avatar    string `json:"avatar"`
	Gid       uint   `json:"gid"`
}

type LoginResponse struct {
	Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	UserInfo UserInfo `json:"user_info"`
}
type UserListResponse struct {
	Response
	Page     int        `json:"page"`
	Total    int        `json:"total"`
	UserList []UserInfo `json:"user_list"`
}
type GroupUser struct {
	Response
	UserList []UserInfo `json:"groupuser_list"`
}
