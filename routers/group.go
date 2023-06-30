package routers

import (
	"Kube-CC/api/v1/group"
	"Kube-CC/middleware"

	"github.com/gin-gonic/gin"
)

func groupRouter(router *gin.RouterGroup) {
	groupRouters := router.Group("/group")
	{
		groupRouters.GET("/info", group.Info)
		groupRouters.GET("/:page", group.Index)                                     // 浏览组信息
		groupRouters.GET("/delete/:id", middleware.Is3Role(), group.Delete)         // 删除组
		groupRouters.POST("/remove/:id", middleware.Is2Role(), group.Remove)        // 移出用户
		groupRouters.POST("/resetpass/:id", middleware.Is2Role(), group.TransAdmin) // 更改管理员
		groupRouters.POST("/update/:id", middleware.Is2Role(), group.Update)        // 更新组信息
	}

}
