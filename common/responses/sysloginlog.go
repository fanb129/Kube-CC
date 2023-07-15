package responses

type LoginInfo struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Role      uint   `json:"role"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Ipaddr    string `json:"ipaddr"`
}

type LogResponse struct {
	Response
	LogID uint `json:"user_id"`
}

type LogInfoResponse struct {
	Response
	LoginInfo LoginInfo `json:"user_info"`
}
type LogListResponse struct {
	Response
	Page    int         `json:"page"`
	Total   int         `json:"total"`
	LogList []LoginInfo `json:"login_list"`
}
