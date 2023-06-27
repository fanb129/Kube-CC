package linux

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

// Index 获取当前用户下的指定类型的linux
func Index(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "")
	uid := 0
	var err error
	if u_id != "" {
		uid, err = strconv.Atoi(u_id)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}
	}
	os := c.DefaultQuery("os", "")
	os1, err := strconv.Atoi(os)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	linuxListResponse, err := service.GetLinux(uint(uid), uint(os1))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, linuxListResponse)
}

// Add 创建指定类型的linux
func Add(c *gin.Context) {
	form := forms.LinuxAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.CreateLinux(form.Uid, form.Kind, form.ExpiredTime, form.Cpu, form.Memory)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Delete 删除linux
func Delete(c *gin.Context) {
	ns := c.Param("name")
	response, err := service.DeleteLinux(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// BatchAdd 批量添加
func BatchAdd(c *gin.Context) {
	// 表单验证
	form := forms.BatchLinuxAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	ids := form.Uid
	group := sync.WaitGroup{}
	group.Add(len(ids))
	for _, id := range ids {
		go func(id uint) {
			if _, err := service.CreateLinux(id, form.Kind, form.ExpiredTime, form.Cpu, form.Memory); err != nil {
				zap.S().Errorln(err)
			}
			group.Done()
		}(id)
	}
	group.Wait()
	c.JSON(http.StatusOK, responses.OK)
}

func Update(c *gin.Context) {
	// 表单验证
	form := forms.LinuxUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	uid := ""
	if form.Uid != 0 {
		uid = strconv.Itoa(int(form.Uid))
	}
	response, err := service.UpdateLinux(form.Name, uid, form.ExpiredTime, form.Cpu, form.Memory)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
