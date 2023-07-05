package pvc

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Index 展示所有pvc
func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	pvcListRsp, err := service.ListPVC(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, pvcListRsp)
	}
}

// Delete 删除指定pvc
func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := service.DeletePVC(ns, name)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// Add 通过表单提交添加pvc
func Add(c *gin.Context) {
	form := forms.PvcAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := service.CreatePVC(form.Namespace, form.Name, form.StorageClassName, form.StorageSize, form.AccessModes)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func Update(c *gin.Context) {
	form := forms.PvcUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	err := service.UpdatePVC(form.Namespace, form.Name, form.StorageSize)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, responses.OK)
}
