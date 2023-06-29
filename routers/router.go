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

	// yaml路由
	yamlRouter(auth)

	//用户路由
	userRouter(auth)

	// node路由
	nodeRouter(auth)

	// namespace路由
	nsRouter(auth)

	// deploy路由
	deployRouter(auth)

	// statefulSet路由
	statefulSetRouter(auth)

	// service路由
	serviceRouter(auth)

	// pod路由
	podRouter(auth)

	// spark路由
	sparkRouter(auth)

	// hadoop路由
	hadoopRouter(auth)

	// linux路由
	linuxRouter(auth)

	return r
}
