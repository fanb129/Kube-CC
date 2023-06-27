package routers

import (
	"Kube-CC/api/v1/spark"
	"Kube-CC/middleware"
	"github.com/gin-gonic/gin"
)

func sparkRouter(router *gin.RouterGroup) {
	sparkRouters := router.Group("/spark")
	{
		//sparkRouter.POST("/add", middleware.Is2Role(), spark.Add) // 新建spark集群
		// 批量添加
		sparkRouters.POST("/add", middleware.Is2Role(), spark.BatchAdd) // 新建spark集群
		sparkRouters.GET("/delete/:name", spark.Delete)
		sparkRouters.GET("", spark.Index)
		sparkRouters.POST("/update", middleware.Is2Role(), spark.Update)
	}
}
