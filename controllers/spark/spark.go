package spark

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
	"strconv"
)

// Index 获取当前用户spark列表
func Index(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "")
	uid, err := strconv.Atoi(u_id)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	sparkListRes, err := service.GetSpark(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, sparkListRes)
}

// Add 创建spark
func Add(c *gin.Context) {
	// 表单验证
	form := common.SparkAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	response, err := service.CreateSpark(form.Uid, form.MasterReplicas, form.WorkerReplicas)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func Delete(c *gin.Context) {
	ns := c.Param("name")
	response, err := service.DeleteSpark(ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
