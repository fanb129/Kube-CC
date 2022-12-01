package deploy

import (
	"Kube-CC/common"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"net/http"
)

// Index 展示所有deploy
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
	deployListResponse, err := service.GetDeploy(ns, selector)
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

// Info 获取单个deploy的yaml信息
func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	res, err := service.GetADeploy(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.DeployInfo{
			Response: common.OK,
			Info:     *res,
		})
	}
}

// Add 通过表单提交添加deploy
//func Add(c *gin.Context) {
//	v1.Deployment{}
//	yamlApply.DeployCreate()
//}
