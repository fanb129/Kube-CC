package responses

// <<新增>>
type GroupInfo struct {
	Groupid     uint   `json:"groupid"`
	CreatedAt   string `json:"group_created_at"`
	UpdatedAt   string `json:"group_updated_at"`
	Name        string `json:"name"`
	Adminid     uint   `json:"adminid"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
}

type GroupInfoResponse struct {
	Response
	GroupInfo GroupInfo `json:"group_info"`
}

//type GroupListResponse struct {
//	Response
//	Page      int         `json:"page"`
//	Total     int         `json:"total"`
//	GroupList []GroupInfo `json:"group_list"`
//}

type GroupListResponse struct {
	Response
	Length    int         `json:"length"`
	GroupList []GroupInfo `json:"group_list"`
}
