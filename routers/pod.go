package routers

import (
	"Kube-CC/api/v1/pod"
	"github.com/gin-gonic/gin"
)

func podRouter(router *gin.RouterGroup) {
	podRouters := router.Group("/pod")
	{
		podRouters.GET("", pod.Index) // 浏览指定namespace的pod
		podRouters.GET("/delete", pod.Delete)
		podRouters.GET("/info", pod.Info)
	}
}
