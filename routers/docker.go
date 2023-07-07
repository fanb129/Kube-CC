package routers

import (
	"Kube-CC/api/v2/docker"
	"github.com/gin-gonic/gin"
)

func dockerRouter(router *gin.RouterGroup) {
	dockerRouters := router.Group("/docker")
	{
		dockerRouters.GET("/index/:page", docker.Index)
		dockerRouters.GET("/remove", docker.Remove)
		dockerRouters.POST("/pullpublic", docker.PullPublic)
		dockerRouters.POST("/tagadd", docker.TagAdd)
		dockerRouters.POST("/pullprivate", docker.PullPrivate)
		dockerRouters.POST("/createimagebyid", docker.CreateImageByImageId)
	}
}
