package dao

import (
	"gorm.io/gorm"
	"k8s_deploy_gin/models"
)

// GetSparkListById 根据uid获取spark列表
func GetSparkListById(u_id uint) ([]models.Spark, error) {
	var sparks []models.Spark
	result := mysqlDb.Where("uid = ?", u_id).Find(&sparks)
	if result.Error != nil {
		return nil, result.Error
	}
	return sparks, nil
}

// CreateSpark 根据uid和t创建spark
func CreateSpark(u_id uint, t string) (int, error) {
	spark := models.Spark{
		Uid:  u_id,
		Time: t,
	}
	rs := mysqlDb.Create(&spark)
	return int(rs.RowsAffected), rs.Error
}

// DeleteSpark 根据spark主键删除
func DeleteSpark(s_id uint) (int, error) {
	spark := models.Spark{
		Model: gorm.Model{ID: s_id},
	}
	rs := mysqlDb.Delete(&spark)
	return int(rs.RowsAffected), rs.Error
}

// GetSpark 根据spark主键获取spark
func GetSpark(s_id uint) (*models.Spark, error) {
	spark := models.Spark{
		Model: gorm.Model{ID: s_id},
	}
	rs := mysqlDb.First(&spark)

	return &spark, rs.Error
}
