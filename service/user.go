package service

import (
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
)

// IndexUser  分页浏览用户信息
func IndexUser(page int) *common.UserListResponse {
	total, userList := dao.GetUserList(page, conf.PageSize)

	// 如果无数据，则返回到第一页
	if total == 0 && page > 1 {
		page = 1
		total, userList = dao.GetUserList(page, conf.PageSize)
	}

	return &common.UserListResponse{
		Response: common.OK,
		Page:     page,
		UserList: userList,
	}
}

// DeleteUSer  删除用户
func DeleteUSer(id int) *common.Response {
	user, err := dao.GetUserById(id)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	user.Status = 0
	err = dao.UpdateUser(user)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	return &common.OK
}

// EditUser 授权用户
func EditUser(id, status int) *common.Response {
	user, err := dao.GetUserById(id)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	user.Status = status
	err = dao.UpdateUser(user)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	return &common.OK
}

// ResetPassUser 重置密码
func ResetPassUser(id int, password string) *common.Response {
	// 获取用户
	user, err := dao.GetUserById(id)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	// 密码加密
	pwd, err := EncryptionPWD(password)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}
	// 修改密码
	user.Password = pwd
	err = dao.UpdateUser(user)
	if err != nil {
		return &common.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	return &common.OK
}
