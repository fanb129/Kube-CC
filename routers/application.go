package routers

import (
	"Kube-CC/api/v1/hadoop"
	"Kube-CC/api/v1/linux"
	"Kube-CC/api/v1/spark"
	"Kube-CC/api/v2/application/deploy"
	"Kube-CC/middleware"
	"github.com/gin-gonic/gin"
)

func appRouter(router *gin.RouterGroup) {
	appRouters := router.Group("/app")

	sparkRouters := appRouters.Group("/spark")
	{
		//sparkRouter.POST("/add", middleware.Is2Role(), spark.Add) // 新建spark集群
		// 批量添加
		sparkRouters.POST("/add", middleware.Is2Role(), spark.BatchAdd) // 新建spark集群
		sparkRouters.GET("/delete/:name", spark.Delete)
		sparkRouters.GET("", spark.Index)
		sparkRouters.POST("/update", middleware.Is2Role(), spark.Update)
	}

	// hadoop 路由
	hadoopRouters := appRouters.Group("/hadoop")
	{
		//hadoopRouter.POST("/add", middleware.Is2Role(), hadoop.Add)
		// 批量添加
		hadoopRouters.POST("/add", middleware.Is2Role(), hadoop.BatchAdd)
		hadoopRouters.GET("/delete/:name", hadoop.Delete)
		hadoopRouters.GET("", hadoop.Index)
		hadoopRouters.POST("/update", middleware.Is2Role(), hadoop.Update)
	}

	// 云主机路由
	linuxRouters := appRouters.Group("/linux")
	{
		linuxRouters.GET("", linux.Index)
		linuxRouters.GET("delete/:name", linux.Delete)
		//linuxRouters.POST("/add", middleware.Is2Role(), linux.Add)
		// 批量添加
		linuxRouters.POST("/add", middleware.Is2Role(), linux.BatchAdd)
		linuxRouters.POST("/update", middleware.Is2Role(), linux.Update)
	}

	// deploy app 无状态应用 路由
	deployRouters := appRouters.Group("/deploy")
	{
		deployRouters.GET("", deploy.Index)
		deployRouters.GET("/delete", deploy.Delete)
		deployRouters.POST("/add", deploy.Add)
		deployRouters.POST("/update", deploy.Update)
		deployRouters.GET("/info", deploy.Info)
	}

	// statefulSet app 有状态应用 路由

	// Job app 一次性任务 路由
}