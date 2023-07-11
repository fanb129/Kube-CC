package routers

import (
	"Kube-CC/api/v1/node"
	"Kube-CC/middleware"
	"github.com/gin-gonic/gin"
)

func nodeRouter(router *gin.RouterGroup) {
	nodeRouters := router.Group("/node")
	{
		nodeRouters.GET("", middleware.Is3Role(), node.Index) // 浏览所有node
		nodeRouters.GET("/delete/:node", middleware.Is3Role(), node.Delete)
		nodeRouters.POST("/add", middleware.Is3Role(), node.Add)
	}
}
