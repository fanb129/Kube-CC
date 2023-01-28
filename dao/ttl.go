package dao

import (
	"Kube-CC/models"
	"time"
)

func ListTtl() ([]models.Ttl, error) {
	var list []models.Ttl
	result := mysqlDb.Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func GetTtlByNs(ns string) (*models.Ttl, error) {
	var ttl models.Ttl
	if result := mysqlDb.Where("namespace = ?", ns).First(&ttl); result.Error != nil {
		return nil, result.Error
	}
	return &ttl, nil
}

// GetDeletedTtlByNs 根据ns查找软删除的ttl
func GetDeletedTtlByNs(ns string) (*models.Ttl, error) {
	ttl := models.Ttl{}
	result := mysqlDb.Unscoped().Where("namespace = ?", ns).First(&ttl)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ttl, nil
}

func CreateTtl(ns string, expiredTime time.Time) (int, error) {
	ttl := models.Ttl{
		Namespace:   ns,
		ExpiredTime: expiredTime,
	}
	result := mysqlDb.Create(&ttl)
	return int(result.RowsAffected), result.Error
}

func DeleteTtl(ttl *models.Ttl) error {
	result := mysqlDb.Delete(ttl)
	return result.Error
}

// UpdateTtl 更新user,零值不会更新
func UpdateTtl(t *models.Ttl) (int, error) {
	result := mysqlDb.Model(t).Updates(models.Ttl{
		ExpiredTime: t.ExpiredTime,
	})
	return int(result.RowsAffected), result.Error
}

func UpdateUnscopedTtl(t *models.Ttl) (int, error) {
	rs := mysqlDb.Unscoped().Save(t)
	return int(rs.RowsAffected), rs.Error
}
