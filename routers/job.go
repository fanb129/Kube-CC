package routers

import (
	"Kube-CC/api/v2/job"
	"github.com/gin-gonic/gin"
)

func jobRouter(router *gin.RouterGroup) {
	podRouters := router.Group("/job")
	{
		podRouters.GET("", job.Index)
		podRouters.GET("/delete", job.Delete)
	}
}
