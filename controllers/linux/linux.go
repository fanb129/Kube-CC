package linux

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
	"strconv"
)

// Index 获取当前用户下的指定类型的linux
func Index(c *gin.Context) {
	u_id, ok := c.Get("u_id")
	kind, _ := strconv.Atoi(c.Param("kind"))
	if !ok {
		c.JSON(http.StatusOK, common.NoUid)
		return
	}
	linuxListResponse, err := service.GetLinux(u_id.(uint), uint(kind))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, linuxListResponse)
}

// Add 创建指定类型的linux
func Add(c *gin.Context) {
	form := common.LinuxAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	response, err := service.CreateLinux(form.Uid, form.Kind)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Delete 删除linux
func Delete(c *gin.Context) {
	ns := c.Param("name")
	response, err := service.DeleteLinux(ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
