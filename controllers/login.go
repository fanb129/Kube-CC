package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"k8s_deploy_gin/dao"
	"time"
)

// Login 用户登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	user := dao.GetUserByName(username)
	// 用户名不存在
	if user == nil {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "username not exists",
			"data": "",
		})
		return
	}
	password := c.PostForm("password")
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// 密码错误
	if err != nil {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "password error",
			"data": "",
		})
	} else {
		// 登陆成功
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "login success",
			"data": user.ToMap(),
		})
	}
}

// Logout 用户登出
func Logout(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "lotout"})
}

// Register 用户注册
func Register(c *gin.Context) {
	//fmt.Println("register")
	data := make(map[string]interface{})
	data["username"] = c.PostForm("username")
	data["nickname"] = c.PostForm("nickname")
	data["password"] = c.PostForm("password")

	user := dao.GetUserByName(data["username"].(string))

	if user != nil {
		// 如果当前用户已注册
		if user.Status > 0 {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "username exists",
				"data": "",
			})
		} else {
			user.Nickname = data["nickname"].(string)
			user.Avatar = ""
			user.Status = 1
			user.Major = ""
			user.Class = ""
			user.CreateAt = time.Now()
			user.Password = EncryptionPWD(data["password"].(string))
			if rs := dao.UpdateUser(user); rs {
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  "register success",
					"data": "",
				})
			} else {
				c.JSON(200, gin.H{
					"code": 0,
					"msg":  "register fail",
					"data": "",
				})
			}
		}
		return
	}

	// 当前用户名不存在时
	data["password"] = EncryptionPWD(data["password"].(string))
	if rs := dao.CreateUser(data); rs {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "register success",
			"data": "",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "register fail",
			"data": "",
		})
	}

}

// EncryptionPWD 对密码进行加密
func EncryptionPWD(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}
