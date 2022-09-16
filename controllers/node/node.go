package node

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"k8s_deploy_gin/service/ws"
	"net/http"
)

func Index(c *gin.Context) {
	nodeListResponse, err := service.GetNode("")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nodeListResponse)
	}
}

func WsSsh(c *gin.Context) {
	fmt.Println("node WsSsh")
	ws.NodeWsSsh(c.Writer, c.Request)
}
