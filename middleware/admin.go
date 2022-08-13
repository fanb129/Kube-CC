package middleware

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"net/http"
)

func Is2Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("role")
		role := r.(uint)
		if role < 2 {
			c.JSON(http.StatusForbidden, common.NoRole)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Is3Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("role")
		role := r.(uint)
		if role < 3 {
			c.JSON(http.StatusForbidden, common.NoRole)
			c.Abort()
			return
		}
		c.Next()
	}
}
