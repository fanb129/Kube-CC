package service

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/conf"
	"Kube-CC/dao"
	"errors"
)

// GetPage  分页浏览日志信息
func GetPage(page int) (*responses.LogListResponse, error) {
	u, total, err := dao.GetLogList(page, conf.PageSize)
	if err != nil {
		return nil, errors.New("获取登录日志失败")
	}
	// 如果无数据，则返回到第一页
	if len(u) == 0 && page > 1 {
		page = 1
		u, total, err = dao.GetLogList(page, conf.PageSize)
		if err != nil {
			return nil, errors.New("获取登录日志失败")
		}
	}
	logList := make([]responses.LoginInfo, len(u))
	for i, v := range u {
		tmp := responses.LoginInfo{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:  v.Username,
		}
		logList[i] = tmp
	}
	return &responses.LogListResponse{
		Response: responses.OK,
		Page:     page,
		Total:    total,
		LogList:  logList,
	}, nil
}

func LogInfo(u_id uint) (*responses.LogInfoResponse, error) {
	log, err := dao.GetLogById(u_id)
	if err != nil {
		return nil, errors.New("获取日志失败")
	}
	return &responses.LogInfoResponse{
		Response: responses.OK,
		LoginInfo: responses.LoginInfo{
			ID:        log.ID,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: log.UpdatedAt.Format("2006-01-02 15:04:05"),
			Username:  log.Username,
			Status:    log.Status,
			Ipaddr:    log.Ipaddr,
		},
	}, nil
}

// DeleteUSer  删除用户
func DeleteLog(id uint) (*responses.Response, error) {
	row, err := dao.DeleteUserById(id)
	if err != nil || row == 0 {
		return nil, errors.New("删除失败")
	}
	return &responses.OK, nil
}

// UpdateLogin 更新日志信息
func UpdateLogin(id uint, data forms.UpdateForm) {

}
