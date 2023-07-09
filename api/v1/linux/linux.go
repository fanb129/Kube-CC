package linux

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/application"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Index 获取当前用户下的指定类型的linux
func Index(c *gin.Context) {
	os := c.Query("os")
	os1, err := strconv.Atoi(os)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	ns := c.Query("ns")
	response, err := application.ListLinux(ns, uint(os1))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// Add 创建指定类型的linux
func Add(c *gin.Context) {
	form := forms.LinuxAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := application.CreateLinux(form.Name, form.Namespace, form.Kind, form.ApplyResources)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Delete 删除linux
func Delete(c *gin.Context) {
	name := c.Query("name")
	ns := c.Query("ns")
	response, err := application.DeleteLinux(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// BatchAdd 批量添加
//func BatchAdd(c *gin.Context) {
//	// 表单验证
//	form := forms.BatchLinuxAddForm{}
//	if err := c.ShouldBind(&form); err != nil {
//		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
//		return
//	}
//	ids := form.Uid
//	group := sync.WaitGroup{}
//	group.Add(len(ids))
//	for _, id := range ids {
//		go func(id uint) {
//			if _, err := application.CreateLinux(id, form.Kind, form.ExpiredTime, form.Resources); err != nil {
//				zap.S().Errorln(err)
//			}
//			group.Done()
//		}(id)
//	}
//	group.Wait()
//	c.JSON(http.StatusOK, responses.OK)
//}

func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.GetLinux(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Update(c *gin.Context) {
	// 表单验证
	form := forms.LinuxUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := application.UpdateLinux(form.Name, form.Namespace, form.ApplyResources)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
