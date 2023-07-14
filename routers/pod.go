package routers

import (
	"Kube-CC/api/v1/pod"
	"github.com/gin-gonic/gin"
)

func podRouter(router *gin.RouterGroup) {
	pvcRouters := router.Group("/pod")
	{
		pvcRouters.GET("/log", pod.Log) // pod 日志
	}
}
