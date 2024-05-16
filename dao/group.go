package dao

import (
	"Kube-CC/models"

	"gorm.io/gorm"
)

//<<新增>>

// GetGroupUserList 分页返回组用户列表(page第几页,pageSize每页几条数据)
//func GetGroupUserList(page int, pageSize int, groupid uint) ([]models.User, int, error) {
//	var users []models.User
//	var total int64
//	mysqlDb.Find(&users).Count(&total)
//	// 计算偏移量 Offset指定开始返回记录前要跳过的记录数。
//	offset := (page - 1) * pageSize
//	// 查看所有的user
//	result := mysqlDb.Where("groupid = ?", groupid).Offset(offset).Limit(pageSize).Find(&users)
//
//	if result.Error != nil {
//		return nil, 0, result.Error
//	}
//	//r := 0
//	//if int(total)%pageSize != 0 {
//	//	r = 1
//	//}
//	//return users, int(total)/pageSize + r, nil
//	return users, int(total), nil
//}
//
//// GetGroupList 分页返回组列表(page第几页,pageSize每页几条数据)
//func GetGroupList(page int, pageSize int) ([]models.Group, int, error) {
//	var groups []models.Group
//	var total int64
//	mysqlDb.Find(&groups).Count(&total)
//	// 计算偏移量 Offset指定开始返回记录前要跳过的记录数。
//	offset := (page - 1) * pageSize
//	// 查看所有的user
//	result := mysqlDb.Offset(offset).Limit(pageSize).Find(&groups)
//
//	if result.Error != nil {
//		return nil, 0, result.Error
//	}
//	//r := 0
//	//if int(total)%pageSize != 0 {
//	//	r = 1
//	//}
//	//return users, int(total)/pageSize + r, nil
//	return groups, int(total), nil
//}

// GetOkUserList 获取可以加入本组的用户
func GetOkUserList() ([]models.User, error) {
	users := []models.User{}
	result := mysqlDb.Where("role = 1 and groupid = 0").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetAllGroup 超级管理元可查看所有的组
func GetAllGroup() ([]models.Group, error) {
	var groups []models.Group
	result := mysqlDb.Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	return groups, nil
}

// GetGroupById 通过id获取group
func GetGroupById(id uint) (*models.Group, error) {
	group := models.Group{}
	result := mysqlDb.First(&group, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}

// GetGroupByAdminid 通过adminid获取Group
func GetGroupByAdminid(adid uint) ([]models.Group, error) {
	groups := []models.Group{}
	result := mysqlDb.Where("adminid = ?", adid).Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	return groups, nil
}

// GetGroupUserById 通过id获取groupuser
func GetGroupUserById(id uint) ([]models.User, error) {
	users := []models.User{}
	result := mysqlDb.Where("groupid = ?", id).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetDeletedGroupByName 根据groupname查找软删除的group
func GetDeletedGroupByName(name string) (*models.Group, error) {
	group := models.Group{}
	result := mysqlDb.Unscoped().Where("name = ?", name).First(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}

// GetGroupByName 通过name获取group
func GetGroupByName(name string) (*models.Group, error) {
	group := models.Group{}
	result := mysqlDb.Where("name = ?", name).First(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}

// DeleteGroupById 根据id删除group
func DeleteGroupById(id uint) (int, error) {
	group := models.Group{
		Model: gorm.Model{
			ID: id,
		},
	}
	result := mysqlDb.Delete(&group)
	return int(result.RowsAffected), result.Error
}

// CreateGroup 新增group
func CreateGroup(adminid uint, name, description string) (int, error) {
	group := models.Group{
		Adminid:     adminid,
		Name:        name,
		Description: description,
	}
	result := mysqlDb.Create(&group)
	return int(result.RowsAffected), result.Error
}
func UpdateUnscopedGroup(g *models.Group) (int, error) {
	rs := mysqlDb.Unscoped().Save(g)
	return int(rs.RowsAffected), rs.Error
}

// UpdateGroupWithNil 更新group,包括零值
func UpdateGroupWithNil(g *models.Group) (int, error) {
	result := mysqlDb.Save(g)
	// result := mysqlDb.Model(g).Updates(models.Group{
	// 	Name: g.Name,
	// 	AdminId: g.AdminId,
	// 	Description: g.Description,
	// })
	return int(result.RowsAffected), result.Error
}

// UpdateUser 更新Group,零值不会更新
func UpdateGroup(g *models.Group) (int, error) {
	result := mysqlDb.Model(g).Updates(models.Group{
		Name:        g.Name,
		Adminid:     g.Adminid,
		Description: g.Description,
	})
	return int(result.RowsAffected), result.Error
}
