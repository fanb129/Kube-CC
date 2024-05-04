package routers

import (
	"Kube-CC/api/v2/docker"
	"github.com/gin-gonic/gin"
)

func dockerRouter(router *gin.RouterGroup) {
	dockerRouters := router.Group("/docker")
	{
		dockerRouters.GET("/", docker.Index)
		dockerRouters.GET("/ok", docker.IndexOk)
		dockerRouters.GET("/delete/:id", docker.Delete)

		dockerRouters.POST("/pull", docker.Pull)
		dockerRouters.POST("/update", docker.Update)
		dockerRouters.POST("/save", docker.Save)
	}
}
