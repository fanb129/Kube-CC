package dao

import (
	"k8s_deploy_gin/models"
	"time"
)

var users []models.User

// GetUserList 分页返回用户列表(page第几页,pageSize每页几条数据)
func GetUserList(page int, pageSize int) (int, []interface{}) {
	// 分页用户列表数据
	userList := make([]interface{}, 0, len(users))
	// 计算偏移量 Offset指定开始返回记录前要跳过的记录数。
	offset := (page - 1) * pageSize
	// 查看所有的user,并获取user总数
	var total int64
	result := MysqlDb.Order("status DESC").Offset(offset).Limit(pageSize).Where("status > ?", "0").Find(&users).Count(&total)
	// 如果返回的数据为0条
	if result.RowsAffected == 0 {
		return 0, userList
	}
	// 返回的user总数,用count更好
	//total := count

	// 查询数据
	result.Offset(offset).Limit(pageSize).Find(&users)
	for _, userSingle := range users {
		userItemMap := userSingle.ToMap()
		userList = append(userList, userItemMap)
	}
	return int(total), userList
}

// GetUserById 通过id获取user
func GetUserById(id int) (*models.User, error) {
	user := models.User{}
	result := MysqlDb.First(&user, id)
	if result.Error == nil {
		return &user, nil
	}
	return nil, result.Error
}

// GetUserByName 通过name获取user
func GetUserByName(name string) (*models.User, error) {
	user := models.User{}
	result := MysqlDb.Where("username = ?", name).First(&user)
	if result.Error == nil {
		return &user, nil
	}
	return nil, result.Error
}

// CreateUser 新增user
func CreateUser(username, nickname, password string) (uint, error) {
	user := models.User{
		Model: models.Model{
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		},
		Username: username,
		Password: password,
		Nickname: nickname,
		Status:   1,
	}
	result := MysqlDb.Create(user)
	if result.Error == nil {
		return user.ID, nil
	}
	return 0, result.Error
}

// UpdateUser 更新user
func UpdateUser(u *models.User) error {
	u.UpdateAt = time.Now()
	result := MysqlDb.Save(&u)
	return result.Error
}
