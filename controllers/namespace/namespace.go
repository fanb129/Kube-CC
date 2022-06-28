package namespace

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
)

// Index 展示所有namespace，即自定义的spark或者hadoop集群
func Index(c *gin.Context) {
	nsListResponse, err := service.GetNs()
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nsListResponse)
	}
}
