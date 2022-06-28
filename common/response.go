package common

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type LoginResponse struct {
	Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}
type RegisterResponse struct {
	Response
	UserID uint `json:"user_id"`
}

var OK = Response{StatusCode: 0, StatusMsg: "success"}
var NoStatus = Response{StatusCode: -1, StatusMsg: "权限不够"}
var NoToken = Response{StatusCode: ErrNoToken.Code, StatusMsg: "No Token"}
var TokenExpired = Response{StatusCode: ErrTokenExpired.Code, StatusMsg: ErrTokenExpired.Message}
var InvalidToken = Response{StatusCode: ErrTokenInvalid.Code, StatusMsg: ErrTokenInvalid.Message}
var InvalidParma = Response{StatusCode: ErrValidateFail.Code, StatusMsg: ErrValidateFail.Message}

type UserListResponse struct {
	Response
	Page     int           `json:"page"`
	UserList []interface{} `json:"user_list"`
}

type NodeListResponse struct {
	Response
	Length   int           `json:"length"`
	NodeList []interface{} `json:"node_list"`
}
type NsListResponse struct {
	Response
	Length int           `json:"length"`
	NsList []interface{} `json:"ns_list"`
}
type PodListResponse struct {
	Response
	Length  int           `json:"length"`
	PodList []interface{} `json:"pod_list"`
}
