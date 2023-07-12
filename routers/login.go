package routers

import (
	"Kube-CC/api/v1/login"

	"github.com/gin-gonic/gin"
)

func loginRouter(router *gin.RouterGroup) {
	router.POST("/login", login.Login)          // 登陆路由
	router.POST("/checkPass", login.CheckPass)  //验证密码
	router.GET("/logout", login.Logout)         // 登出路由
	router.POST("/register", login.Register)    // 注册路由
	router.GET("/captcha", login.GetCaptcha)    // 验证码
	router.POST("/checkcp", login.CheckCaptcha) // 验证验证码
}
