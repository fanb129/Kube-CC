package dao

import (
	"Kube-CC/models"
	"errors"
	"gorm.io/gorm"
)

func GetImageById(id uint) (*models.Docker, error) {
	docker := models.Docker{}
	result := mysqlDb.First(&docker, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &docker, nil
}

// GetAllImages 获取所有用户的所有状态image
func GetAllImages() ([]models.Docker, error) {
	var dockers []models.Docker
	result := mysqlDb.Find(&dockers)
	if result.Error != nil {
		return nil, result.Error
	}
	return dockers, nil
}

// GetPrivateImage 看自己的私有
func GetPrivateImage(uid uint) ([]models.Docker, error) {
	var dockers []models.Docker
	result := mysqlDb.Where("user_id = ? and kind = 2", uid).Find(&dockers)
	if result.Error != nil {
		return nil, result.Error
	}
	return dockers, nil
}

// GetPublicImage 获取共有
func GetPublicImage() ([]models.Docker, error) {
	var dockers []models.Docker
	result := mysqlDb.Where("kind = 1").Find(&dockers)
	if result.Error != nil {
		return nil, result.Error
	}
	return dockers, nil
}

// GetImagesByUid 用户能够看到的镜像，包括自己的私有，和别人的共有
func GetImagesByUid(uid uint) ([]models.Docker, error) {
	dockers, err := GetPrivateImage(uid)
	if err != nil {
		return nil, err
	}
	dockers1, err := GetPublicImage()
	if err != nil {
		return nil, err
	}

	dockers = append(dockers, dockers1...)
	return dockers, nil
}

// GetOkIMages 应用部署时的可用image，包括public和自己的，但必须是status==1
func GetOkIMages(uid uint) ([]models.Docker, error) {
	var dockers, dockers1 []models.Docker
	// 自己的私有
	result := mysqlDb.Where("user_id = ? and kind = 2 and status = 1", uid).Find(&dockers)
	if result.Error != nil {
		return nil, result.Error
	}

	// 公有
	result = mysqlDb.Where("kind = 1 and status = 1").Find(&dockers1)
	if result.Error != nil {
		return nil, result.Error
	}

	dockers = append(dockers, dockers1...)
	return dockers, nil
}

func CreateImage(name, tag, description string, uid, kind uint) (uint, error) {
	var docker models.Docker
	result := mysqlDb.Where("image_name = ? and tag = ? and user_id = ?", name, tag, uid).First(&docker)
	if result.RowsAffected > 0 {
		return 0, errors.New("镜像已存在")
	}
	docker.ImageName = name
	docker.Tag = tag
	docker.UserId = uid
	docker.Kind = kind
	docker.Status = 2
	docker.Description = description

	result = mysqlDb.Create(&docker)
	return docker.ID, result.Error
}

func SaveImage(d *models.Docker) error {
	result := mysqlDb.Save(d)
	return result.Error
}

func UpdateImage(id uint, kind uint, description string) error {
	docker := models.Docker{
		Model: gorm.Model{
			ID: id,
		},
	}
	result := mysqlDb.Model(&docker).Updates(models.Docker{
		Kind:        kind,
		Description: description,
	})
	return result.Error
}

func DeleteImage(id uint) error {
	docker := models.Docker{
		Model: gorm.Model{
			ID: id,
		},
	}
	result := mysqlDb.Delete(&docker)
	return result.Error
}

func DeleteUserAllPrivateImages(uid uint) error {
	result := mysqlDb.Where("user_id = ? and kind = 2", uid).Delete(&models.Docker{})
	return result.Error
}
