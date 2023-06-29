package svc

import (
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index 展示所有service
func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	//u_id := c.DefaultQuery("u_id", "")
	selector := ""
	//if u_id != "" {
	//	label := map[string]string{
	//		"u_id": u_id,
	//	}
	//	// 将map标签转换为string
	//	selector = labels.SelectorFromSet(label).String()
	//}
	serviceListResponse, err := service.GetService(ns, selector)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, serviceListResponse)
	}
}

func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := service.DeleteService(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	res, err := service.GetAService(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, responses.ServiceInfo{
			Response: responses.OK,
			Info:     *res,
		})
	}
}

//func Add(c *gin.Context) {
//	v1.Service{}
//}
