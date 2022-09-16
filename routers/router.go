package routers

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/controllers"
	"k8s_deploy_gin/controllers/deploy"
	"k8s_deploy_gin/controllers/hadoop"
	"k8s_deploy_gin/controllers/linux"
	"k8s_deploy_gin/controllers/namespace"
	"k8s_deploy_gin/controllers/node"
	"k8s_deploy_gin/controllers/pod"
	"k8s_deploy_gin/controllers/spark"
	"k8s_deploy_gin/controllers/svc"
	"k8s_deploy_gin/controllers/user"
	"k8s_deploy_gin/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())             // 日志
	r.Use(middleware.CorsHandler()) // 跨域设置
	r.Use(gin.Recovery())           // 恐慌 恢复
	gin.SetMode(conf.AppMode)

	r.GET("/api/node/ssh", node.WsSsh)
	r.GET("/api/pod/ssh", pod.Ssh)
	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/login", controllers.Login) // 登陆路由
		apiRouter.POST("/checkPass", controllers.CheckPass)
		apiRouter.GET("/logout", controllers.Logout)      // 登出路由
		apiRouter.POST("/register", controllers.Register) // 注册路由
	}

	// 需要鉴权
	auth := apiRouter.Group("", middleware.JWTToken())
	//用户路由
	userRouter := auth.Group("/user")
	{
		userRouter.GET("/info", user.Info)
		userRouter.GET("/:page", user.Index)                                    // 浏览用户信息
		userRouter.GET("/delete/:id", middleware.Is2Role(), user.Delete)        // 删除用户
		userRouter.POST("/edit/:id", middleware.Is2Role(), user.Edit)           // 授权用户
		userRouter.POST("/resetpass/:id", middleware.Is2Role(), user.ResetPass) // 重置密码
		userRouter.POST("/update/:id", user.Update)                             // 更新用户信息
	}

	// node路由
	nodeRouter := auth.Group("/node")
	{
		nodeRouter.GET("", node.Index) // 浏览所有node
	}

	// namespace路由
	nsRouter := auth.Group("/ns")
	{
		nsRouter.GET("", namespace.Index)                          // 浏览所有namespace
		nsRouter.GET("/delete/:ns", namespace.Delete)              // 删除指定namespace
		nsRouter.POST("/add", middleware.Is3Role(), namespace.Add) // 添加namespace
	}

	// deploy路由
	deployRouter := auth.Group("/deploy")
	{
		deployRouter.GET("", deploy.Index)
		deployRouter.GET("/delete", deploy.Delete)
	}

	serviceRouter := auth.Group("/service")
	{
		serviceRouter.GET("", svc.Index)
		serviceRouter.GET("/delete", svc.Delete)
	}
	// pod路由
	podRouter := auth.Group("/pod")
	{
		podRouter.GET("", pod.Index) // 浏览指定namespace的pod
		podRouter.GET("/delete", pod.Delete)
	}

	// spark路由
	sparkRouter := auth.Group("/spark")
	{
		sparkRouter.POST("/add", middleware.Is2Role(), spark.Add) // 新建spark集群
		sparkRouter.GET("/delete/:name", spark.Delete)
		sparkRouter.GET("", spark.Index)
	}

	hadoopRouter := auth.Group("/hadoop")
	{
		hadoopRouter.POST("/add", middleware.Is2Role(), hadoop.Add)
		hadoopRouter.GET("/delete/:name", hadoop.Delete)
		hadoopRouter.GET("", hadoop.Index)
	}

	// linux路由
	linuxRouter := auth.Group("/linux")
	{
		linuxRouter.GET("", linux.Index)
		linuxRouter.GET("delete/:name", linux.Delete)
		linuxRouter.POST("/add", middleware.Is2Role(), linux.Add)
	}

	return r
}
