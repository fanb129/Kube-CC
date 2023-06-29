package routers

import (
	"Kube-CC/api/v2/docker"
	"github.com/gin-gonic/gin"
)

func dockerRouter(router *gin.RouterGroup) {
	dockerRouters := router.Group("/docker")
	{
		dockerRouters.GET("/info", docker.Info)
		dockerRouters.GET("/:page", docker.Index)

		//TODO 待定补充镜像操作的一些权限管理
		dockerRouters.GET("/remove/:id", docker.Remove)
		dockerRouters.POST("/update/:id", docker.Update)
		dockerRouters.POST("/save/:id", docker.Save)
		dockerRouters.POST("/pull/:id", docker.Pull)
	}
}
