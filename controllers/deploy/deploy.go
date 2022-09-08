package deploy

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
)

// Index 展示所有deploy
func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	deployListResponse, err := service.GetDeploy(ns, "")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, deployListResponse)
	}
}

// Delete 删除指定deploy
func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := service.DeleteDeploy(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
