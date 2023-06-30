package responses

import "Kube-CC/common"

type Response struct {
	StatusCode int    `json:"code"`
	StatusMsg  string `json:"msg,omitempty"`
}

type ResponseOfValidator struct {
	StatusCode int         `json:"code"`
	StatusMsg  interface{} `json:"msg,omitempty"`
}

func ValidatorResponse(err error) ResponseOfValidator {
	return ResponseOfValidator{
		-1,
		common.Translate(err),
	}
}

var OK = Response{StatusCode: 1, StatusMsg: "success"}
var NoRole = Response{StatusCode: -1, StatusMsg: "权限不够"}
var NoToken = Response{StatusCode: 50008, StatusMsg: "No Token"}
var TokenExpired = Response{StatusCode: 50014, StatusMsg: "token过期"}
var NoUid = Response{StatusCode: -1, StatusMsg: "Uid获取失败"}
var NoGid = Response{StatusCode: -1, StatusMsg: "Gid获取失败"}
