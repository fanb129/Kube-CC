package routers

import (
	"Kube-CC/api/v1/yaml"

	"github.com/gin-gonic/gin"
)

func yamlRouter(router *gin.RouterGroup) {
	yamlRouters := router.Group("/yaml")
	{
		yamlRouters.POST("/apply", yaml.Apply)
		yamlRouters.POST("/create", yaml.Create)
	}
}
