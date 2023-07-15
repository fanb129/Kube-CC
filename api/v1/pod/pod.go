package pod

import (
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"Kube-CC/service/ws/podSsh"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ssh(c *gin.Context) {
	fmt.Println("pod ssh")
	podSsh.PodWsSsh(c.Writer, c.Request)
}

func Log(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	log, err := service.GetPodLog(ns, name)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, log)
}
