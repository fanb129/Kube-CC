package routers

import (
	"Kube-CC/api/v1/deploy"
	"Kube-CC/api/v1/hadoop"
	"Kube-CC/api/v1/linux"
	"Kube-CC/api/v1/login"
	"Kube-CC/api/v1/namespace"
	"Kube-CC/api/v1/node"
	"Kube-CC/api/v1/pod"
	"Kube-CC/api/v1/spark"
	"Kube-CC/api/v1/svc"
	"Kube-CC/api/v1/user"
	"Kube-CC/api/v1/yaml"
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
	{
		apiRouter.POST("/login", login.Login) // 登陆路由
		apiRouter.POST("/checkPass", login.CheckPass)
		apiRouter.GET("/logout", login.Logout)      // 登出路由
		apiRouter.POST("/register", login.Register) // 注册路由
		apiRouter.GET("/Captcha", login.GetCaptcha) // 验证码
	}

	// 需要鉴权
	auth := apiRouter.Group("", middleware.JWTToken())

	// yaml路由
	yamlRouter := auth.Group("/yaml")
	{
		yamlRouter.POST("/apply", yaml.Apply)
		yamlRouter.POST("/create", yaml.Create)
	}

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
		nodeRouter.GET("/delete/:node", middleware.Is3Role(), node.Delete)
		nodeRouter.POST("/add", middleware.Is3Role(), node.Add)
	}

	// namespace路由
	nsRouter := auth.Group("/ns")
	{
		nsRouter.GET("", namespace.Index)                                // 浏览所有namespace
		nsRouter.GET("/delete/:ns", namespace.Delete)                    // 删除指定namespace
		nsRouter.POST("/add", middleware.Is2Role(), namespace.Add)       // 添加namespace
		nsRouter.POST("/update", middleware.Is2Role(), namespace.Update) // 更新uid
	}

	// deploy路由
	deployRouter := auth.Group("/deploy")
	{
		deployRouter.GET("", deploy.Index)
		deployRouter.GET("/delete", deploy.Delete)
		deployRouter.GET("/info", deploy.Info)
		//deployRouter.POST("/add", deploy.Add) // 通过表单添加deploy
	}

	serviceRouter := auth.Group("/service")
	{
		serviceRouter.GET("", svc.Index)
		serviceRouter.GET("/delete", svc.Delete)
		serviceRouter.GET("/info", svc.Info)
	}
	// pod路由
	podRouter := auth.Group("/pod")
	{
		podRouter.GET("", pod.Index) // 浏览指定namespace的pod
		podRouter.GET("/delete", pod.Delete)
		podRouter.GET("/info", pod.Info)
	}

	// spark路由
	sparkRouter := auth.Group("/spark")
	{
		//sparkRouter.POST("/add", middleware.Is2Role(), spark.Add) // 新建spark集群
		// 批量添加
		sparkRouter.POST("/add", middleware.Is2Role(), spark.BatchAdd) // 新建spark集群
		sparkRouter.GET("/delete/:name", spark.Delete)
		sparkRouter.GET("", spark.Index)
		sparkRouter.POST("/update", middleware.Is2Role(), spark.Update)
	}

	hadoopRouter := auth.Group("/hadoop")
	{
		//hadoopRouter.POST("/add", middleware.Is2Role(), hadoop.Add)
		// 批量添加
		hadoopRouter.POST("/add", middleware.Is2Role(), hadoop.BatchAdd)
		hadoopRouter.GET("/delete/:name", hadoop.Delete)
		hadoopRouter.GET("", hadoop.Index)
		hadoopRouter.POST("/update", middleware.Is2Role(), hadoop.Update)
	}

	// linux路由
	linuxRouter := auth.Group("/linux")
	{
		linuxRouter.GET("", linux.Index)
		linuxRouter.GET("delete/:name", linux.Delete)
		//linuxRouter.POST("/add", middleware.Is2Role(), linux.Add)
		// 批量添加
		linuxRouter.POST("/add", middleware.Is2Role(), linux.BatchAdd)
		linuxRouter.POST("/update", middleware.Is2Role(), linux.Update)
	}

	return r
}
