package dao

import (
	"Kube-CC/models"
	"time"

	"gorm.io/gorm"
)

func GetLogList(page int, pageSize int) ([]models.SysLoginLog, int, error) {
	var logs []models.SysLoginLog
	var total int64
	mysqlDb.Find(&logs).Count(&total)

	offset := (page - 1) * pageSize

	result := mysqlDb.Offset(offset).Limit(pageSize).Find(&logs)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return logs, int(total), nil
}

// GetUserById 通过id获取log
func GetLogById(id uint) (*models.SysLoginLog, error) {
	logs := models.SysLoginLog{}
	result := mysqlDb.First(&logs, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &logs, nil
}

// GetDeletedLogByName 根据username查找软删除的log
func GetDeletedLogByName(name string) (*models.SysLoginLog, error) {
	logs := models.SysLoginLog{}
	result := mysqlDb.Unscoped().Where("username = ?", name).First(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return &logs, nil
}

// GetLogByName 通过name获取log
func GetLogByName(name string) (*models.SysLoginLog, error) {
	logs := models.SysLoginLog{}
	result := mysqlDb.Where("username = ?", name).First(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return &logs, nil
}

// DeleteLogById 根据id删除log
func DeleteLogById(id uint) (int, error) {
	logs := models.SysLoginLog{
		Model: gorm.Model{
			ID: id,
		},
	}
	result := mysqlDb.Delete(&logs)
	return int(result.RowsAffected), result.Error
}

// CreateLog 新增log  <<修改>>
func CreateLog(username, status string, logintime time.Time, ipaddr, loginlocation, browser, os string) (int, error) {
	logs := models.SysLoginLog{
		Username:      username,
		Status:        status,
		LoginTime:     logintime,
		Ipaddr:        ipaddr,
		LoginLocation: loginlocation,
		Browser:       browser,
		Os:            os,
	}
	result := mysqlDb.Create(&logs)
	return int(result.RowsAffected), result.Error
}

func UpdateUnscopedLogin(u *models.SysLoginLog) (int, error) {
	rs := mysqlDb.Unscoped().Save(u)
	return int(rs.RowsAffected), rs.Error
}

// UpdateUserWithNil 更新user,包括零值
func UpdateLoginWithNil(u *models.SysLoginLog) (int, error) {
	result := mysqlDb.Save(u)
	return int(result.RowsAffected), result.Error
}

// UpdateUser 更新user,零值不会更新 <<修改>>
func UpdateLogin(u *models.SysLoginLog) (int, error) {
	result := mysqlDb.Model(u).Updates(models.SysLoginLog{
		Username:      u.Username,
		Status:        u.Status,
		LoginTime:     u.LoginTime,
		Ipaddr:        u.Ipaddr,
		LoginLocation: u.LoginLocation,
		Browser:       u.Browser,
		Os:            u.Os,
	})
	return int(result.RowsAffected), result.Error
}

type GeneralDelDto struct {
	Id  int   `uri:"id" json:"id" validate:"required"`
	Ids []int `json:"ids"`
}

func (g GeneralDelDto) GetIds() []int {
	ids := make([]int, 0)
	if g.Id != 0 {
		ids = append(ids, g.Id)
	}
	if len(g.Ids) > 0 {
		for _, id := range g.Ids {
			if id > 0 {
				ids = append(ids, id)
			}
		}
	} else {
		if g.Id > 0 {
			ids = append(ids, g.Id)
		}
	}
	if len(ids) <= 0 {
		//方式全部删除
		ids = append(ids, 0)
	}
	return ids
}
