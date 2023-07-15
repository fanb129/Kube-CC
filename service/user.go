package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"errors"
	"k8s.io/apimachinery/pkg/labels"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// IndexUser  分页浏览用户信息
func IndexUser(page int) (*responses.UserListResponse, error) {
	u, total, err := dao.GetUserList(page, conf.PageSize)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}
	// 如果无数据，则返回到第一页
	if len(u) == 0 && page > 1 {
		page = 1
		u, total, err = dao.GetUserList(page, conf.PageSize)
		if err != nil {
			return nil, errors.New("获取用户列表失败")
		}
	}
	userList := make([]responses.UserInfo, len(u))
	for i, v := range u {
		tmp := responses.UserInfo{
			ID:          v.ID,
			CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:    v.Username,
			Nickname:    v.Nickname,
			Role:        v.Role,
			Avatar:      v.Avatar,
			Gid:         v.Groupid,
			Cpu:         v.Cpu,
			Memory:      v.Memory,
			Storage:     v.Storage,
			PvcStorage:  v.Pvcstorage,
			Gpu:         v.Gpu,
			ExpiredTime: v.ExpiredTime.Format("2006-01-02 15:04:05"),
		}
		userList[i] = tmp
	}
	return &responses.UserListResponse{
		Response: responses.OK,
		Page:     page,
		Total:    total,
		UserList: userList,
	}, nil
}

func GetAll() (*responses.AllUserList, error) {
	u, err := dao.GetAllUser()
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}
	userList := make([]responses.UserInfo, len(u))
	for i, v := range u {
		tmp := responses.UserInfo{
			ID:          v.ID,
			CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:    v.Username,
			Nickname:    v.Nickname,
			Role:        v.Role,
			Avatar:      v.Avatar,
			Gid:         v.Groupid,
			Cpu:         v.Cpu,
			Memory:      v.Memory,
			Storage:     v.Storage,
			PvcStorage:  v.Pvcstorage,
			Gpu:         v.Gpu,
			ExpiredTime: v.ExpiredTime.Format("2006-01-02 15:04:05"),
		}
		userList[i] = tmp
	}
	return &responses.AllUserList{
		Response: responses.OK,
		UserList: userList,
	}, nil
}

func UserInfo(u_id uint) (*responses.UserInfoResponse, error) {
	user, err := dao.GetUserById(u_id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	return &responses.UserInfoResponse{
		Response: responses.OK,
		UserInfo: responses.UserInfo{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:    user.Username,
			Nickname:    user.Nickname,
			Role:        user.Role,
			Avatar:      user.Avatar,
			Gid:         user.Groupid,
			Cpu:         user.Cpu,
			Memory:      user.Memory,
			Storage:     user.Storage,
			PvcStorage:  user.Pvcstorage,
			Gpu:         user.Gpu,
			ExpiredTime: user.ExpiredTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

// // GetAdminUser 获取管理员用户
// func GetAdminUser() (*responses.AdminUserListResponse, error) {
// 	u, err := dao.GetAdmin()
// 	if err != nil {
// 		return nil, errors.New("获取管理员用户失败")
// 	}
// 	adminuserList := make([]responses.AdminUserResponse, len(u))
// 	for i, v := range u {
// 		tmp := responses.AdminUserResponse{
// 			ID:       v.ID,
// 			Username: v.Username,
// 		}
// 		adminuserList[i] = tmp
// 	}
// 	return &responses.AdminUserListResponse{
// 		Response:      responses.OK,
// 		AdminUserList: adminuserList,
// 	}, nil
// }

// DeleteUSer  删除用户
func DeleteUSer(id uint) (*responses.Response, error) {
	row, err := dao.DeleteUserById(id)
	if err != nil || row == 0 {
		return nil, errors.New("删除失败")
	}
	// 删除其所有ns
	label := map[string]string{
		"u_id": strconv.Itoa(int(id)),
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	nsList, err := ListNs(selector)
	if err != nil {
		return nil, err
	}
	for _, ns := range nsList.NsList {
		DeleteNs(ns.Name)
	}
	return &responses.OK, nil
}

// EditUser 授权用户
func EditUser(id, role uint) (*responses.Response, error) {
	user, err := dao.GetUserById(id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	user.Role = role
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("更新失败")
	}

	return &responses.OK, nil
}

func AllocationUser(id uint, data forms.AllocationForm) (*responses.Response, error) {
	user, err := dao.GetUserById(id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	errcpu := VerifyCpu(data.Cpu)
	if errcpu != nil {
		return nil, errors.New("Cpu配额格式错误")
	}
	errmemory := VerifyResource(data.Memory)
	if errmemory != nil {
		return nil, errors.New("内存配额格式错误")
	}
	errstorage := VerifyResource(data.Storage)
	if errstorage != nil {
		return nil, errors.New("存储配额格式错误")
	}
	errpvcstorage := VerifyResource(data.Pvcstorage)
	if errpvcstorage != nil {
		return nil, errors.New("持久化存储配额格式错误")
	}
	errgpu := VerifyResource(data.Gpu)
	if errgpu != nil {
		return nil, errors.New("Gpu配额格式错误")
	}
	var timeLayoutStr = "2006-01-02 15:04:05"
	fmexptime, errtime := time.ParseInLocation(timeLayoutStr, data.ExpiredTime, time.Local)
	if errtime != nil {
		return nil, errors.New("日期转换错误")
	}
	user.Cpu = data.Cpu
	user.Memory = data.Memory
	user.Storage = data.Storage
	user.Pvcstorage = data.Pvcstorage
	user.Gpu = data.Gpu
	user.ExpiredTime = fmexptime
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("更新用户配额失败")
	}
	return &responses.OK, nil
}

// UpdateUser 更新用户信息
func UpdateUser(id uint, data forms.UpdateForm) (*responses.Response, error) {
	user, err := dao.GetUserById(id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	user.Avatar = data.Avatar
	user.Nickname = data.Nickname
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("更新失败")
	}

	return &responses.OK, nil
}

// ResetPassUser 重置密码
func ResetPassUser(id uint, password string) (*responses.Response, error) {
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

	return &responses.OK, nil
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
