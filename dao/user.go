package dao

import (
	"golang.org/x/crypto/bcrypt"
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
	result := db.Order("status DESC").Offset(offset).Limit(pageSize).Where("status > ?", "0").Find(&users).Count(&total)
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
func GetUserById(id int) *models.User {
	user := models.User{}
	db.First(&user, id)
	return &user
}

// GetUserByName 通过name获取user
func GetUserByName(name string) *models.User {
	user := models.User{}
	result := db.Where("username = ?", name).First(&user)
	if result.RowsAffected == 0 {
		return nil
	}
	return &user
}

// CreateUser 新增user
func CreateUser(data map[string]interface{}) bool {
	result := db.Create(&models.User{
		Model: models.Model{
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		},
		Username: data["username"].(string),
		Password: data["password"].(string),
		Nickname: data["nickname"].(string),
		Status:   1,
	})

	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// UpdateUser 更新user
func UpdateUser(u *models.User) bool {
	u.UpdateAt = time.Now()
	result := db.Save(&u)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// CheckUser 检查密码是否正确
func CheckUser(username string, password string) bool {
	user := GetUserByName(username)
	if user == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
