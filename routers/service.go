package routers

import (
	"Kube-CC/api/v1/svc"
	"github.com/gin-gonic/gin"
)

func serviceRouter(router *gin.RouterGroup) {
	serviceRouters := router.Group("/service")
	{
		serviceRouters.GET("", svc.Index)
		serviceRouters.GET("/delete", svc.Delete)
		serviceRouters.GET("/info", svc.Info)
	}
}
