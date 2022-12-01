package spark

import (
	"Kube-CC/common"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

// Index 获取当前用户spark列表
func Index(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "")
	uid := 0
	var err error
	if u_id != "" {
		uid, err = strconv.Atoi(u_id)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}
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

func Update(c *gin.Context) {
	// 表单验证
	form := common.SparkUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	uid := ""
	if form.Uid != 0 {
		uid = strconv.Itoa(int(form.Uid))
	}
	response, err := service.UpdateSpark(form.Name, uid, form.MasterReplicas, form.WorkerReplicas)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// BatchAdd 批量添加
func BatchAdd(c *gin.Context) {
	// 表单验证
	form := common.BatchSparkAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	ids := form.Uid
	group := sync.WaitGroup{}
	group.Add(len(ids))
	for _, id := range ids {
		go func(id uint) {
			if _, err := service.CreateSpark(id, form.MasterReplicas, form.WorkerReplicas); err != nil {
				zap.S().Errorln(err)
			}
			group.Done()
		}(id)
	}
	group.Wait()
	c.JSON(http.StatusOK, common.OK)
}
