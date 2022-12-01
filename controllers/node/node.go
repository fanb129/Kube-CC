package node

import (
	"Kube-CC/common"
	"Kube-CC/service"
	"Kube-CC/service/ssh"
	"Kube-CC/service/ws"
	"fmt"
	"github.com/gin-gonic/gin"
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

func Add(c *gin.Context) {
	form := common.NodeAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	nodes := make([]ssh.Config, len(form.Nodes))
	for i, node := range form.Nodes {
		tmp := ssh.Config{
			Host:     node.Host,
			Port:     form.Port,
			User:     form.User,
			Password: form.Password,
			Type:     ssh.TypePassword,
		}
		nodes[i] = tmp
	}
	rsp, err := service.CreateNode(nodes)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	node := c.Param("name")
	rsp, err := service.DeleteNode(node)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rsp)
}
