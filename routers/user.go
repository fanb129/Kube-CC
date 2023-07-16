package routers

import (
	"Kube-CC/api/v1/user"
	"Kube-CC/middleware"

	"github.com/gin-gonic/gin"
)

func userRouter(router *gin.RouterGroup) {
	userRouters := router.Group("/user")
	{
		userRouters.GET("/info", user.Info)
		userRouters.GET("/:page", user.Index)                             // 浏览用户信息
		userRouters.GET("/delete/:id", middleware.Is2Role(), user.Delete) // 删除用户
		userRouters.GET("/getall", middleware.Is2Role(), user.GetAll)
		userRouters.POST("/edit/:id", middleware.Is2Role(), user.Edit)             // 授权用户
		userRouters.POST("/resetpass/:id", middleware.Is2Role(), user.ResetPass)   // 重置密码
		userRouters.POST("/update/:id", user.Update)                               // 更新用户信息
		userRouters.POST("/allocation/:id", middleware.Is2Role(), user.Allocation) // 修改用户配额
	}

}
