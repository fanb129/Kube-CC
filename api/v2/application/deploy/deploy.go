package deploy

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/application"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Add(c *gin.Context) {
	form := forms.DeployAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := application.CreateAppDeploy(form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.DeleteAppDeploy(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	if ns == "" {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: "请选择namespace"})
		return
	}
	response, err := application.ListAppDeploy(ns, "")
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Update(c *gin.Context) {
	form := forms.DeployAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := application.UpdateAppDeploy(form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.GetAppDeploy(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
