package pod

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
)

func GetPod(c *gin.Context) {
	ns := c.DefaultPostForm("namespace", "default")
	podListResponse, err := service.GetPod(ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, podListResponse)
	}
}
