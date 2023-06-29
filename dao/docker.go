package dao

import (
	"Kube-CC/models"
)

// 根据用户获取用户对应的镜像列表
func GetImageList(page int, pageSize int) ([]models.Docker, int, error) {
	var dockers []models.Docker
	var total int64
	mysqlDb.Find(&dockers).Count(&total)

	offset := (page - 1) * pageSize

	result := mysqlDb.Offset(offset).Limit(pageSize).Find(&dockers)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return dockers, int(total), result.Error
}

func GetImgById(imageid string) (*models.Docker, error) {
	imglist := models.Docker{}
	result := mysqlDb.Where("ImageId = ?", imageid).Find(&imglist)
	if result.Error != nil {
		return nil, result.Error
	}
	return &imglist, nil
}

// 根据镜像id删除镜像

func GetDeletedImgById(id string) (*models.Docker, error) {
	imglist := models.Docker{}
	result := mysqlDb.Unscoped().Where("ImageId = ?", id).First(&imglist)
	if result.Error != nil {
		return nil, result.Error
	}
	return &imglist, nil
}

// 新增镜像

func CreateImage(imageid string, uid uint, kinds int) (int, error) {
	img := models.Docker{
		ImageId: imageid,
		UserId:  uid,
		Kind:    kinds,
	}
	result := mysqlDb.Create(&img)
	return int(result.RowsAffected), result.Error
}

// 更新tag
func UpdateImage(u *models.Docker) (int, error) {
	result := mysqlDb.Model(u).Updates(models.Docker{
		ImageId: u.ImageId,
	})
	return int(result.RowsAffected), result.Error
}
