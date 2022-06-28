package models

type User struct {
	Model
	Username string `json:"username"` // 账号
	Password string `json:"password"` // 密码
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Status   int    `json:"status"`   // 权限状态 0删除 1普通用户 2集群管理员 3管理员(上课老师) 4超级管理员
	Major    string `json:"major"`    // 专业
	Class    string `json:"class"`    // 班级
}

// TableName 设置User的表名为user
func (User) TableName() string {
	return "user"
}

func (u User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       u.ID,
		"username": u.Username,
		"nickname": u.Nickname,
		"avatar":   u.Avatar,
		"status":   u.Status,
		"major":    u.Major,
		"class":    u.Class,
		"createAt": u.CreateAt,
		"updateAt": u.UpdateAt,
	}
}
