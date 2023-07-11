package dao

import (
	"Kube-CC/conf"
	"Kube-CC/models"

	"gorm.io/gorm"

	"regexp"
)

// GetUserList 分页返回用户列表(page第几页,pageSize每页几条数据)
func GetUserList(page int, pageSize int) ([]models.User, int, error) {
	var users []models.User
	var total int64
	mysqlDb.Find(&users).Count(&total)
	// 计算偏移量 Offset指定开始返回记录前要跳过的记录数。
	offset := (page - 1) * pageSize
	// 查看所有的user
	result := mysqlDb.Offset(offset).Limit(pageSize).Find(&users)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	//r := 0
	//if int(total)%pageSize != 0 {
	//	r = 1
	//}
	//return users, int(total)/pageSize + r, nil
	return users, int(total), nil
}

// GetUserById 通过id获取user
func GetUserById(id uint) (*models.User, error) {
	user := models.User{}
	result := mysqlDb.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetDeletedUserByName 根据username查找软删除的user
func GetDeletedUserByName(name string) (*models.User, error) {
	user := models.User{}
	result := mysqlDb.Unscoped().Where("username = ?", name).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetDeletedUserByEmail 根据email查找软删除的user <<新增>>
func GetDeletedUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	result := mysqlDb.Unscoped().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByName 通过name获取user
func GetUserByName(name string) (*models.User, error) {
	user := models.User{}
	result := mysqlDb.Where("username = ?", name).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail 通过Email获取user <<新增>>
func GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	result := mysqlDb.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// // GetAdmin 获取所有管理员用户
// func GetAdmin() ([]models.User, error) {
// 	var users []models.User
// 	result := mysqlDb.Where("role = ?", 2).Find(&users)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return users, nil
// }

// DeleteUserById 根据id删除user
func DeleteUserById(id uint) (int, error) {
	user := models.User{
		Model: gorm.Model{
			ID: id,
		},
	}
	result := mysqlDb.Delete(&user)
	return int(result.RowsAffected), result.Error
}

// DeleteUserByEmail 根据email删除user <<新增>>
func DeleteUserByEmail(email string) (int, error) {
	user := models.User{
		Email: email,
	}
	result := mysqlDb.Delete(&user)
	return int(result.RowsAffected), result.Error
}

// CreateUser 新增user  <<修改>>
func CreateUser(username, nickname, password, email string) (int, error) {
	user := models.User{
		Username:   username,
		Email:      email,
		Nickname:   nickname,
		Password:   password,
		Cpu:        conf.Cpu,
		Memory:     conf.Memory,
		Storage:    conf.Storage,
		Pvcstorage: conf.Pvcstorage,
		Gpu:        conf.Gpu,
	}
	result := mysqlDb.Create(&user)
	return int(result.RowsAffected), result.Error
}
func UpdateUnscopedUser(u *models.User) (int, error) {
	rs := mysqlDb.Unscoped().Save(u)
	return int(rs.RowsAffected), rs.Error
}

// UpdateUserWithNil 更新user,包括零值
func UpdateUserWithNil(u *models.User) (int, error) {
	result := mysqlDb.Save(u)
	//result := mysqlDb.Model(u).Updates(map[string]interface{}{
	//	"created_at": u.CreatedAt,
	//	"deleted_at": gorm.DeletedAt{},
	//	"username":   u.Username,
	//	"password":   u.Password,
	//	"nickname":   u.Nickname,
	//	"role":       u.Role,
	//	"avatar":     u.Avatar,
	//})
	return int(result.RowsAffected), result.Error
}

// UpdateUser 更新user,零值不会更新 <<修改>>
func UpdateUser(u *models.User) (int, error) {
	result := mysqlDb.Model(u).Updates(models.User{
		Username:   u.Username,
		Email:      u.Email,
		Nickname:   u.Nickname,
		Password:   u.Password,
		Role:       u.Role,
		Avatar:     u.Avatar,
		Groupid:    u.Groupid,
		Cpu:        u.Cpu,
		Memory:     u.Memory,
		Storage:    u.Storage,
		Pvcstorage: u.Pvcstorage,
		Gpu:        u.Gpu,
	})
	return int(result.RowsAffected), result.Error
}

// 匹配邮箱 <<新增>>
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
