package job

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/application"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Add(c *gin.Context) {
	form := forms.JobAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := application.CreateAppJob(form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, rsp)
}

func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	if ns == "" {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: "请选择namespace"})
		return
	}
	response, err := application.ListAppJob(ns, "")
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.DeleteAppJob(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.GetAppJob(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
