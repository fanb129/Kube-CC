package node

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"Kube-CC/service/ws"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	nodeListResponse, err := service.GetNode("")
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nodeListResponse)
	}
}

func WsSsh(c *gin.Context) {
	fmt.Println("node WsSsh")
	ws.NodeWsSsh(c.Writer, c.Request)
}

//func Add(c *gin.Context) {
//	form := forms.NodeAddForm{}
//	if err := c.ShouldBind(&form); err != nil {
//		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
//		return
//	}
//	nodes := make([]ssh.Config, len(form.Nodes))
//	for i, node := range form.Nodes {
//		tmp := ssh.Config{
//			Host:     node.Host,
//			Port:     form.Port,
//			User:     form.User,
//			Password: form.Password,
//			Type:     ssh.TypePassword,
//		}
//		nodes[i] = tmp
//	}
//	rsp, err := service.CreateNode(nodes)
//	if err != nil {
//		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, rsp)
//}

func Add(c *gin.Context) {
	form := forms.NodeAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	hosts := make([]string, len(form.Nodes))
	for i, node := range form.Nodes {
		hosts[i] = node.Host
	}
	// 异步处理
	go func() {
		rsp, err := service.CreateNodeBySealos(form.Password, hosts)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, rsp)
	}()
	rsp := responses.Response{StatusCode: 1, StatusMsg: "操作成功,请稍后刷新"}
	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	node := c.Param("node")
	// 异步删除
	go func() {
		rsp, err := service.DeleteNodeBysealos(node)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		c.JSON(http.StatusOK, rsp)
	}()
	rsp := responses.Response{StatusCode: 1, StatusMsg: "操作成功，请稍后刷新"}
	c.JSON(http.StatusOK, rsp)
}
