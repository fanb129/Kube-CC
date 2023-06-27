package routers

import (
	"Kube-CC/api/v1/linux"
	"Kube-CC/middleware"
	"github.com/gin-gonic/gin"
)

func linuxRouter(router *gin.RouterGroup) {
	linuxRouters := router.Group("/linux")
	{
		linuxRouters.GET("", linux.Index)
		linuxRouters.GET("delete/:name", linux.Delete)
		//linuxRouters.POST("/add", middleware.Is2Role(), linux.Add)
		// 批量添加
		linuxRouters.POST("/add", middleware.Is2Role(), linux.BatchAdd)
		linuxRouters.POST("/update", middleware.Is2Role(), linux.Update)
	}
}
