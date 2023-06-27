package middleware

import (
	"Kube-CC/common/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Is2Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("role")
		role := r.(uint)
		if role < 2 {
			c.JSON(http.StatusForbidden, responses.NoRole)
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
			c.JSON(http.StatusForbidden, responses.NoRole)
			c.Abort()
			return
		}
		c.Next()
	}
}
