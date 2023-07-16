package dao

import (
	"Kube-CC/models"
)

// 根据用户获取用户对应的镜像列表
func GetImageList(page int, pageSize int, uid uint) ([]models.Docker, int, error) {
	var images []models.Docker
	var total int64
	// uid为0时显示全部镜像
	if uid == 0 {
		mysqlDb.Find(&images).Count(&total)
		offset := (page - 1) * pageSize
		// 当查询失败时返回0
		if result := mysqlDb.Offset(offset).Limit(pageSize).Find(&images); result.Error != nil {
			return nil, 0, result.Error
		} else {
			return images, int(total), result.Error
		}
	} else {
		mysqlDb.Find(&images).Count(&total)
		offset := (page - 1) * pageSize
		// 当查询失败时返回0
		if result := mysqlDb.Offset(offset).Limit(pageSize).Where("user_id = ?", uid).Find(&images); result.Error != nil {
			return nil, 0, result.Error
		} else {
			return images, int(total), result.Error
		}
	}

	return images, 0, nil
}

func GetImgById(imageid string) (*models.Docker, error) {
	img := models.Docker{}
	result := mysqlDb.First(&img, imageid)
	if result.Error != nil {
		return nil, result.Error
	}
	return &img, nil
}

func GetImgByName(imagename string) (*models.Docker, error) {
	img := models.Docker{}
	result := mysqlDb.Where("image_name = ?", imagename).First(&img)
	if result.Error != nil {
		return nil, result.Error
	}
	return &img, nil
}

// 根据镜像id删除镜像

func DeletedImgByImageId(id string) (*models.Docker, error) {
	imglist := models.Docker{}
	result := mysqlDb.Where("image_id = ?", id).Delete(&models.Docker{})
	if result.Error != nil {
		return nil, result.Error
	}
	return &imglist, nil
}

// 新增镜像

func CreateImage(imagename, imageid string, uid uint, kinds int, tag, size string) (int, error) {
	img := models.Docker{
		ImageName: imagename,
		ImageId:   imageid,
		UserId:    uid,
		Kind:      kinds,
		Tag:       tag,
		Size:      size,
	}
	result := mysqlDb.Create(&img)
	return int(result.RowsAffected), result.Error
}

// 更新tag
func CreateImageByTag(u *models.Docker) (int, error) {
	result := mysqlDb.Model(u).Updates(models.Docker{
		Tag: u.Tag,
	})
	return int(result.RowsAffected), result.Error
}
