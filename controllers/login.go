package controllers

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
)

// Login 用户登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	login, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, login)
	}
}

// Logout 用户登出
func Logout(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "lotout"})
}

// Register 用户注册
func Register(c *gin.Context) {
	//fmt.Println("register")
	username := c.PostForm("username")
	nickname := c.PostForm("nickname")
	password := c.PostForm("password")
	register, err := service.Register(username, nickname, password)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, register)
	}
}
