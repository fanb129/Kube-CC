package routers

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/controllers"
	"k8s_deploy_gin/controllers/namespace"
	"k8s_deploy_gin/controllers/node"
	"k8s_deploy_gin/controllers/pod"
	"k8s_deploy_gin/controllers/spark"
	"k8s_deploy_gin/controllers/user"
	"k8s_deploy_gin/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())             // 日志
	r.Use(middleware.CorsHandler()) // 跨域设置
	r.Use(gin.Recovery())           // 恐慌 恢复
	gin.SetMode(conf.AppMode)

	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/login", controllers.Login)       // 登陆路由
		apiRouter.GET("/logout", controllers.Logout)      // 登出路由
		apiRouter.POST("/register", controllers.Register) // 注册路由
	}

	// 需要鉴权
	auth := apiRouter.Group("", middleware.JWTToken())
	//用户路由
	userRouter := auth.Group("/user")
	{
		userRouter.GET("/:page", user.Index)                                    // 浏览用户信息
		userRouter.GET("/delete/:id", middleware.Is3Role(), user.Delete)        // 删除用户
		userRouter.POST("/edit/:id", middleware.Is3Role(), user.Edit)           // 授权用户
		userRouter.POST("/resetpass/:id", middleware.Is3Role(), user.ResetPass) // 重置密码
	}

	// node路由
	nodeRouter := auth.Group("/node")
	{
		nodeRouter.GET("", node.Index) // 浏览所有node
	}

	// namespace路由
	nsRouter := auth.Group("/ns")
	{
		nsRouter.GET("", namespace.Index)           // 浏览所有namespace
		nsRouter.POST("/delete/", namespace.Delete) // 删除指定namespace
		nsRouter.POST("/add", namespace.Add)        // 添加namespace
	}

	// pod路由
	podRouter := auth.Group("/pod")
	{
		podRouter.POST("", pod.GetPod) // 浏览指定namespace的pod
	}

	// spark路由
	sparkRouter := auth.Group("/spark")
	{
		sparkRouter.POST("/add", middleware.Is2Role(), spark.Add) // 新建spark集群
		sparkRouter.GET("/delete/:s_id", spark.Delete)
		sparkRouter.GET("", spark.Index)
	}
	return r
}
