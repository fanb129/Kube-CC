package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"errors"
)

// <<新增>>
// IndexGroup  分页浏览组信息
func IndexGroup(page int) (*responses.GroupListResponse, error) {
	g, total, err := dao.GetGroupList(page, conf.PageSize)
	if err != nil {
		return nil, errors.New("获取组列表失败")
	}
	// 如果无数据，则返回到第一页
	if len(g) == 0 && page > 1 {
		page = 1
		g, total, err = dao.GetGroupList(page, conf.PageSize)
		if err != nil {
			return nil, errors.New("获取组列表失败")
		}
	}
	groupList := make([]responses.GroupInfo, len(g))
	for i, v := range g {
		tmp := responses.GroupInfo{
			Groupid:     v.ID,
			CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Name:        v.Name,
			Adminid:     v.Adminid,
			Description: v.Description,
		}
		groupList[i] = tmp
	}
	return &responses.GroupListResponse{
		Response:  responses.OK,
		Page:      page,
		Total:     total,
		GroupList: groupList,
	}, nil
}

// IndexGroupUser 分页浏览组内信息
func IndexGroupUser(page int, groupid uint) (*responses.UserListResponse, error) {
	gu, total, err := dao.GetGroupUserList(page, conf.PageSize, groupid)
	if err != nil {
		return nil, errors.New("获取该组用户列表失败")
	}
	// 如果无数据，则返回到第一页
	if len(gu) == 0 && page > 1 {
		page = 1
		gu, total, err = dao.GetGroupUserList(page, conf.PageSize, groupid)
		if err != nil {
			return nil, errors.New("获取该组用户列表失败")
		}
	}
	groupuserList := make([]responses.UserInfo, len(gu))
	for i, v := range gu {
		tmp := responses.UserInfo{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:  v.Username,
			Nickname:  v.Nickname,
			Role:      v.Role,
			Avatar:    v.Avatar,
		}
		groupuserList[i] = tmp
	}
	return &responses.UserListResponse{
		Response: responses.OK,
		Page:     page,
		Total:    total,
		UserList: groupuserList,
	}, nil
}

// ViewGroupUser 查看组内成员
func ViewGroupUser(groupid uint) (*responses.GroupUser, error) {
	gu, err := dao.GetGroupUserById(groupid)
	if err != nil || len(gu) == 0 {
		return nil, errors.New("获取该组用户列表失败")
	}
	groupuserList := make([]responses.UserInfo, len(gu))
	for i, v := range gu {
		tmp := responses.UserInfo{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:  v.Username,
			Nickname:  v.Nickname,
			Role:      v.Role,
			Avatar:    v.Avatar,
		}
		groupuserList[i] = tmp
	}
	return &responses.GroupUser{
		Response: responses.OK,
		UserList: groupuserList,
	}, nil
}

func GroupInfo(g_id uint) (*responses.GroupInfoResponse, error) {
	group, err := dao.GetGroupById(g_id)
	if err != nil {
		return nil, errors.New("获取组失败")
	}
	return &responses.GroupInfoResponse{
		Response: responses.OK,
		GroupInfo: responses.GroupInfo{
			Groupid:     group.ID,
			CreatedAt:   group.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   group.UpdatedAt.Format("2006-01-02 15:04:05"),
			Name:        group.Name,
			Adminid:     group.Adminid,
			Description: group.Description,
		},
	}, nil
}

// DeleteGroup  删除组
func DeleteGroup(id uint) (*responses.Response, error) {
	users, erru := dao.GetGroupUserById(id)
	if erru != nil || users == nil {
		return nil, errors.New("获取组用户失败")
	}
	for _, v := range users {
		v.Groupid = 0
	}
	row, err := dao.DeleteGroupById(id)
	if err != nil || row == 0 {
		return nil, errors.New("删除失败")
	}
	return &responses.OK, nil
}

// RemoveUser 从组中移出用户
func RemoveUser(u_id uint) (*responses.Response, error) {
	user, err := dao.GetUserById(u_id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	user.Groupid = 0
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("移出失败")
	}

	return &responses.OK, nil
}

// AddUser 向组中添加用户
func AddUser(g_id, u_id uint) (*responses.Response, error) {
	user, err := dao.GetUserById(u_id)
	if err != nil {
		return nil, errors.New("获取用户失败")
	}
	if user.Groupid != 0 {
		return nil, errors.New("该用户属于其他组")
	}
	user.Groupid = g_id
	row, err := dao.UpdateUser(user)
	if err != nil || row == 0 {
		return nil, errors.New("添加失败")
	}

	return &responses.OK, nil
}

// UpdateGroup 更新组信息
func UpdateGroup(g_id uint, data forms.GroupUpdateForm) (*responses.Response, error) {
	group, err := dao.GetGroupById(g_id)
	if err != nil {
		return nil, errors.New("获取组失败")
	}
	group.Description = data.Description
	group.Name = data.Name
	row, err := dao.UpdateGroup(group)
	if err != nil || row == 0 {
		return nil, errors.New("更新失败")
	}

	return &responses.OK, nil
}

// TransAdmin 转移管理员
func TransAdmin(g_id, odad_id, nwad_id uint) (*responses.Response, error) {
	// 获取组和用户
	group, err := dao.GetGroupById(g_id)
	if err != nil {
		return nil, errors.New("获取组失败")
	}
	_, erro := dao.GetUserById(odad_id)
	if erro != nil {
		return nil, errors.New("获取旧管理员失败")
	}
	if group.Adminid != odad_id {
		return nil, errors.New("旧管理员不是本组管理员")
	}
	nwad, errn := dao.GetUserById(nwad_id)
	if errn != nil {
		return nil, errors.New("获取新管理员失败")
	}
	if nwad.Groupid != g_id {
		return nil, errors.New("新管理员不属于本组")
	}
	group.Adminid = nwad_id

	row, err := dao.UpdateGroup(group)
	if err != nil || row == 0 {
		return nil, errors.New("更新失败")
	}

	return &responses.OK, nil
}
