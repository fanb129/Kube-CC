package responses

// <<新增>>
type GroupInfo struct {
	GroupID     uint   `json:"groupid"`
	CreatedAt   string `json:"group_created_at"`
	UpdatedAt   string `json:"group_updated_at"`
	Name        string `json:"name"`
	AdminID     uint   `json:"adminid"`
	Description string `json:"description"`
}

type GroupInfoResponse struct {
	Response
	GroupInfo GroupInfo `json:"group_info"`
}
type GroupListResponse struct {
	Response
	Page      int         `json:"page"`
	Total     int         `json:"total"`
	GroupList []GroupInfo `json:"group_list"`
}
