package routers

import (
	"Kube-CC/api/v1/pod"
	"github.com/gin-gonic/gin"
)

func podRouter(router *gin.RouterGroup) {
	podRouters := router.Group("/pod")
	{
		podRouters.GET("/log", pod.Log) // pod 日志
		podRouters.GET("/delete", pod.Delete)
	}
}
