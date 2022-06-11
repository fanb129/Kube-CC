package routers

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/controllers"
	"k8s_deploy_gin/controllers/user"
	"k8s_deploy_gin/middleware/cors"
	"k8s_deploy_gin/middleware/myjwt"
	"k8s_deploy_gin/pkg/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())       // 日志
	r.Use(cors.CorsHandler()) // 跨域设置
	r.Use(gin.Recovery())     // 恐慌 恢复
	gin.SetMode(setting.RunMode)
	var authMiddleware = myjwt.GinJWTMiddlewareInit(&myjwt.AllUserAuthorizator{})

	//404 handler
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 0, "msg": " Page not found", "data": ""})
	})

	auth := r.Group("/auth")
	{
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	r.POST("/login", authMiddleware.LoginHandler) // 登陆路由
	r.GET("/logout", controllers.Logout)          // 登出路由
	r.POST("/register", controllers.Register)     // 注册路由

	// api路径下需要登录
	apiRouter := r.Group("/api")
	apiRouter.Use(authMiddleware.MiddlewareFunc())

	fourMiddleware := myjwt.GinJWTMiddlewareInit(&myjwt.FourAuthorizator{})
	//threeMiddleware := myjwt.GinJWTMiddlewareInit(&myjwt.ThreeAuthorizator{})
	//twoMiddleware := myjwt.GinJWTMiddlewareInit(&myjwt.TwoAuthorizator{})

	//用户路由
	userRouter := apiRouter.Group("/user")
	{
		userRouter.GET("/:page", user.Index)                                               // 浏览用户信息
		userRouter.GET("/delete/:id", user.Delete, fourMiddleware.MiddlewareFunc())        // 删除用户
		userRouter.POST("/edit/:id", user.Edit, fourMiddleware.MiddlewareFunc())           // 授权用户
		userRouter.POST("/resetpass/:id", user.ResetPass, fourMiddleware.MiddlewareFunc()) // 重置密码
	}
	return r
}
