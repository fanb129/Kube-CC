package routers

import (
	"Kube-CC/api/v2/statefulSet"
	"github.com/gin-gonic/gin"
)

func statefulSetRouter(router *gin.RouterGroup) {
	statefulSetRouters := router.Group("/deploy")
	{
		statefulSetRouters.GET("", StatefulSet.Index)
		statefulSetRouters.GET("/delete", StatefulSet.Delete)
		statefulSetRouters.GET("/info", StatefulSet.Info)
		//statefulSetRouters.POST("/add", statefulSet.Add) // 通过表单添加statefulSet
	}
}
