package routers

import (
	"Kube-CC/api/v1/hadoop"
	"Kube-CC/api/v1/linux"
	"Kube-CC/api/v1/spark"
	"Kube-CC/api/v2/application/deploy"
	"Kube-CC/api/v2/application/statefulSet"
	"Kube-CC/middleware"
	"github.com/gin-gonic/gin"
)

func appRouter(router *gin.RouterGroup) {
	appRouters := router.Group("/app")

	sparkRouters := appRouters.Group("/spark")
	{
		sparkRouters.POST("/add", spark.Add) // 新建spark集群
		// 批量添加
		sparkRouters.POST("/batchadd", middleware.Is2Role(), spark.BatchAdd) // 新建spark集群
		sparkRouters.GET("/delete/:name", spark.Delete)
		sparkRouters.GET("", spark.Index)
		sparkRouters.POST("/update", spark.Update)
		sparkRouters.GET("/info", spark.Info)
	}

	// hadoop 路由
	hadoopRouters := appRouters.Group("/hadoop")
	{
		hadoopRouters.POST("/add", hadoop.Add)
		// 批量添加
		hadoopRouters.POST("/batchadd", middleware.Is2Role(), hadoop.BatchAdd)
		hadoopRouters.GET("/delete/:name", hadoop.Delete)
		hadoopRouters.GET("", hadoop.Index)
		hadoopRouters.POST("/update", hadoop.Update)
		hadoopRouters.GET("/info", hadoop.Info)
	}

	// 云主机路由
	linuxRouters := appRouters.Group("/linux")
	{
		linuxRouters.GET("", linux.Index)
		linuxRouters.GET("/delete", linux.Delete)
		linuxRouters.POST("/add", linux.Add)
		// 批量添加
		//linuxRouters.POST("/add", middleware.Is2Role(), linux.BatchAdd)
		linuxRouters.POST("/update", linux.Update)
		linuxRouters.GET("/info", linux.Info)
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
	statefulSetRouters := appRouters.Group("/statefulSet")
	{
		statefulSetRouters.GET("", statefulSet.Index)
		statefulSetRouters.GET("/delete", statefulSet.Delete)
		statefulSetRouters.POST("/add", statefulSet.Add)
		statefulSetRouters.POST("/update", statefulSet.Update)
		statefulSetRouters.GET("/info", statefulSet.Info)
	}
	// Job app 一次性任务 路由
}
