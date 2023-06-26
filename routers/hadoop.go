package routers

import (
	"Kube-CC/api/v1/hadoop"
	"Kube-CC/middleware"
	"github.com/gin-gonic/gin"
)

func hadoopRouter(router *gin.RouterGroup) {
	hadoopRouters := router.Group("/hadoop")
	{
		//hadoopRouter.POST("/add", middleware.Is2Role(), hadoop.Add)
		// 批量添加
		hadoopRouters.POST("/add", middleware.Is2Role(), hadoop.BatchAdd)
		hadoopRouters.GET("/delete/:name", hadoop.Delete)
		hadoopRouters.GET("", hadoop.Index)
		hadoopRouters.POST("/update", middleware.Is2Role(), hadoop.Update)
	}
}
