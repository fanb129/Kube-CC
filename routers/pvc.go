package routers

import (
	"Kube-CC/api/v2/pvc"
	"github.com/gin-gonic/gin"
)

func pvcRouter(router *gin.RouterGroup) {
	pvcRouters := router.Group("/pvc")
	{
		pvcRouters.GET("", pvc.Index)          // 浏览所有pvc
		pvcRouters.GET("/delete", pvc.Delete)  // 删除指定pvc
		pvcRouters.POST("/add", pvc.Add)       // 添加pvc
		pvcRouters.POST("/update", pvc.Update) // 更新pvc
	}
}
