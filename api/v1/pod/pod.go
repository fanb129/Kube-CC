package pod

import (
	"Kube-CC/service/ws/podSsh"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Ssh(c *gin.Context) {
	fmt.Println("pod ssh")
	podSsh.PodWsSsh(c.Writer, c.Request)
}
