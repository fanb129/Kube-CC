package routers

import (
	"Kube-CC/api/v1/node"
	"Kube-CC/api/v1/pod"
	"Kube-CC/conf"
	"Kube-CC/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	//r.Use(gin.Logger())             // 日志
	r.Use(middleware.CorsHandler()) // 跨域设置
	r.Use(gin.Recovery())           // 恐慌 恢复
	gin.SetMode(conf.AppMode)

	r.GET("/api/node/ssh", node.WsSsh)
	r.GET("/api/pod/ssh", pod.Ssh)

	apiRouter := r.Group("/api")
	// 登陆路由
	loginRouter(apiRouter)

	// 需要鉴权
	auth := apiRouter.Group("", middleware.JWTToken())

	//用户路由
	userRouter(auth)

	//组路由
	groupRouter(auth)

	// node路由
	nodeRouter(auth)

	// namespace路由
	nsRouter(auth)

	// docker路由
	dockerRouter(auth)

	// pvc持久卷路由
	pvcRouter(auth)

	// sc存储类路由
	scRouter(auth)

	// application路由
	appRouter(auth)

	podRouter(auth)

	return r
}
