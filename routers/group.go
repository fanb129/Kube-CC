package routers

import (
	"Kube-CC/api/v1/group"
	"Kube-CC/middleware"

	"github.com/gin-gonic/gin"
)

func groupRouter(router *gin.RouterGroup) {
	groupRouters := router.Group("/group")
	{
		//groupRouters.GET("/info", group.Info)
		groupRouters.GET("/index", group.Index)                             // 通过当前管理员id
		groupRouters.GET("/okuser", middleware.Is2Role(), group.OkUser)     // 获取可加入组的用户
		groupRouters.GET("/all", middleware.Is3Role(), group.All)           //
		groupRouters.GET("/delete/:id", middleware.Is2Role(), group.Delete) // 删除组
		//groupRouters.GET("/view/:id", middleware.Is2Role(), group.ViewGroupUser) // 查看组成员
		groupRouters.POST("/creat", middleware.Is2Role(), group.Create)      // 创建新的组
		groupRouters.POST("/add/:id", middleware.Is2Role(), group.Add)       // 添加用户
		groupRouters.GET("/remove/:id", middleware.Is2Role(), group.Remove)  // 移出用户
		groupRouters.POST("/update/:id", middleware.Is2Role(), group.Update) // 更新组信息
		//groupRouters.GET("/vgbyad/:id", group.ViewGroupByAdid)                   //通过管理员id查看组
	}

}
