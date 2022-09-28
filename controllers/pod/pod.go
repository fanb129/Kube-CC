package pod

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"k8s_deploy_gin/service/ws/podSsh"
	"net/http"
)

func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	u_id := c.DefaultQuery("u_id", "")
	selector := ""
	if u_id != "" {
		label := map[string]string{
			"u_id": u_id,
		}
		// 将map标签转换为string
		selector = labels.SelectorFromSet(label).String()
	}
	podListResponse, err := service.GetPod(ns, selector)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, podListResponse)
	}
}

func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := service.DeletePod(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Ssh(c *gin.Context) {
	fmt.Println("pod ssh")
	podSsh.PodWsSsh(c.Writer, c.Request)
}

func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	res, err := service.GetAPod(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.PodInfo{
			Response: common.OK,
			Info:     *res,
		})
	}
}
