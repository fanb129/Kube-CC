package myjwt

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/models"
	"k8s_deploy_gin/pkg/setting"
	"log"
	"time"
)

var identityKey = setting.IdentityKey

func GinJWTMiddlewareInit(jwtAuthorizator IAuthorizator) (authMiddleware *jwt.GinJWTMiddleware) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 5,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if u, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"username": u.Username,
					"status":   u.Status,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Username: claims["username"].(string),
				//Status:   claims["status"].(int),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			//handles the login logic. On success LoginResponse is called, on failure Unauthorized is called
			var loginVals models.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if dao.CheckUser(userID, password) {
				return &models.User{
					Username: userID,
					Password: password,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: jwtAuthorizator.HandleAuthorizator,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header:Authorization, query:token, cookie:jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value.
		//This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
