package routers

import (
	"Kube-CC/api/v2/sc"
	"github.com/gin-gonic/gin"
)

func scRouter(router *gin.RouterGroup) {
	scRouters := router.Group("/sc")
	{
		scRouters.GET("", sc.Index) // 浏览所有pvc
	}
}
