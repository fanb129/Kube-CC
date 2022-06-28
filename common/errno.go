package common

type Errno struct {
	Code    int
	Message string
}

var (
	// 数据库相关 101开头
	ErrDataBase            = &Errno{Code: 10101, Message: "数据库错误"}
	ErrQueryUserInfoFail   = &Errno{Code: 10102, Message: "查询用户信息错误"}
	ErrQueryUserLoginFail  = &Errno{Code: 10103, Message: "查询用户登录信息错误"}
	ErrCreateUserFail      = &Errno{Code: 10104, Message: "创建用户信息错误"}
	ErrCreateUserLoginFail = &Errno{Code: 10105, Message: "创建用户登录信息失败"}

	// Token相关 102开头
	ErrTokenExpired   = &Errno{Code: 10201, Message: "Token已过期"}
	ErrTokenSetupFail = &Errno{Code: 10202, Message: "Token生成失败"}
	ErrNoToken        = &Errno{Code: 10203, Message: "No Token"}
	ErrTokenInvalid   = &Errno{Code: 10204, Message: "不是一个Token"}

	// 用户相关 103开头
	ErrPassWordWrong       = &Errno{Code: 10301, Message: "密码错误"}
	ErrEncryptPassWordFail = &Errno{Code: 10302, Message: "密码加密失败"}
	ErrQueryUserNameFail   = &Errno{Code: 10303, Message: "获取用户名失败"}

	// 数据验证相关 104开头

	ErrValidateFail = &Errno{Code: 10401, Message: "数据验证失败"}
)
