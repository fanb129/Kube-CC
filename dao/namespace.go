package dao

import "k8s_deploy_gin/models"

// GetNsList 根据uid和image获取namespace切片
func GetNsList(uid uint, image string) ([]models.Namespace, error) {
	var nsList []models.Namespace
	if uid == 0 {
		// image不为空
		if result := mysqlDb.Where("image = ?", image).Find(&nsList); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if image == "" {
			if result := mysqlDb.Where("uid = ?", uid).Find(&nsList); result.Error != nil {
				return nil, result.Error
			}
		} else {
			if result := mysqlDb.Where("uid = ? AND image = ?", uid, image).Find(&nsList); result.Error != nil {
				return nil, result.Error
			}
		}
	}
	return nsList, nil
}

// GetNsByName 根据ns获取namespace
func GetNsByName(name string) (*models.Namespace, error) {
	ns := models.Namespace{}
	result := mysqlDb.Where("ns = ?", name).First(&ns)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ns, nil
}

// UpdateNsWithNil 更新namespace，包括零值
func UpdateNsWithNil(ns *models.Namespace) (int, error) {
	result := mysqlDb.Save(ns)
	return int(result.RowsAffected), result.Error
}

// DeleteNsByName 根据ns删除namespace
func DeleteNsByName(name string) (int, error) {
	result := mysqlDb.Where("ns = ?", name).Delete(&models.Namespace{})
	return int(result.RowsAffected), result.Error
}

// CreateNs 创建namespace
func CreateNs(uid uint, image, name string) (int, error) {
	ns := models.Namespace{
		Uid:   uid,
		Image: image,
		Ns:    name,
	}
	result := mysqlDb.Create(&ns)
	return int(result.RowsAffected), result.Error
}
