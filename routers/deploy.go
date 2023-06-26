package routers

import (
	"Kube-CC/api/v1/deploy"
	"github.com/gin-gonic/gin"
)

func deployRouter(router *gin.RouterGroup) {
	deployRouters := router.Group("/deploy")
	{
		deployRouters.GET("", deploy.Index)
		deployRouters.GET("/delete", deploy.Delete)
		deployRouters.GET("/info", deploy.Info)
		//deployRouters.POST("/add", deploy.Add) // 通过表单添加deploy
	}
}
