package routers

import (
	"Kube-CC/api/v1/namespace"
	"github.com/gin-gonic/gin"
)

func nsRouter(router *gin.RouterGroup) {
	nsRouters := router.Group("/ns")
	{
		nsRouters.GET("", namespace.Index)                 // 浏览所有namespace
		nsRouters.GET("/delete/:ns", namespace.Delete)     // 删除指定namespace
		nsRouters.POST("/add", namespace.Add)              // 添加namespace
		nsRouters.POST("/update", namespace.Update)        // 更新资源大小
		nsRouters.GET("/total", namespace.NsTotal)         /// 总的资源
		nsRouters.GET("/allkind", namespace.ListAllKindNs) /// 所有类型的ns
	}
}
