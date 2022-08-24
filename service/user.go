package service

import (
	"errors"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
	"math/rand"
	"strings"
	"time"
)

// IndexUser  分页浏览用户信息
func IndexUser(page int) (*common.UserListResponse, error) {
	u, err := dao.GetUserList(page, conf.PageSize)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}
	// 如果无数据，则返回到第一页
	if len(u) == 0 && page > 1 {
		page = 1
		u, err = dao.GetUserList(page, conf.PageSize)
		if err != nil {
			return nil, errors.New("获取用户列表失败")
		}
	}
	userList := make([]common.UserInfo, len(u))
	for i, v := range u {
		tmp := common.UserInfo{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Username:  v.Username,
			Nickname:  v.Nickname,
			Role:      v.Role,
			Avatar:    v.Avatar,
		}
		userList[i] = tmp
	}
	return &common.UserListResponse{
		Response: common.OK,
		Page:     page,
		UserList: userList,
	}, nil
}

// DeleteUSer  删除用户
func DeleteUSer(id uint) (*common.Response, error) {
	row, err := dao.DeleteUserById(id)
	if err != nil || row == 0 {
		return nil, errors.New("删除失败")
	}
	return &common.OK, nil
}

// EditUser 授权用户
func EditUser(id, role uint) (*common.Response, error) {
	user, err := dao.GetUserById(id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	user.Role = role
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("更新失败")
	}

	return &common.OK, nil
}

// ResetPassUser 重置密码
func ResetPassUser(id uint, password string) (*common.Response, error) {
	// 获取用户
	user, err := dao.GetUserById(id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	// 密码加密
	pwd, err := EncryptionPWD(password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	// 修改密码
	user.Password = pwd
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("更新失败")
	}

	return &common.OK, nil
}

func CreatePWD(n int) string {

	pwd := strings.Builder{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		ind := rand.Intn(62)
		if ind >= 0 && ind <= 9 {
			pwd.WriteByte(byte('0' + ind))
		} else if ind <= 35 {
			pwd.WriteByte(byte('a' + ind - 10))
		} else {
			pwd.WriteByte(byte('A' + ind - 36))
		}
	}
	return pwd.String()
}
