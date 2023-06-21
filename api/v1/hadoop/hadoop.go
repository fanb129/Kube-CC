package hadoop

import (
	"Kube-CC/common"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

// Index 获取当前用户的Hadoop列表
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
	hadoopListResponse, err := service.GetHadoop(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, hadoopListResponse)
}

// Add 创建hadoop
func Add(c *gin.Context) {
	// 表单验证
	form := common.HadoopAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}

	res, err := service.CreateHadoop(form.Uid, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas, form.ExpiredTime, form.Cpu, form.Memory)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Delete 根据hadoop名字删除hadoop
func Delete(c *gin.Context) {
	ns := c.Param("name")
	res, err := service.DeleteHadoop(ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func Update(c *gin.Context) {
	// 表单验证
	form := common.HadoopUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	uid := ""
	if form.Uid != 0 {
		uid = strconv.Itoa(int(form.Uid))
	}
	res, err := service.UpdateHadoop(form.Name, uid, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas, form.ExpiredTime, form.Cpu, form.Memory)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// BatchAdd 批量添加
func BatchAdd(c *gin.Context) {
	// 表单验证
	form := common.BatchHadoopAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	ids := form.Uid
	group := sync.WaitGroup{}
	group.Add(len(ids))
	for _, id := range ids {
		go func(id uint) {
			if _, err := service.CreateHadoop(id, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas, form.ExpiredTime, form.Cpu, form.Memory); err != nil {
				zap.S().Errorln(err)
			}
			group.Done()
		}(id)
	}
	group.Wait()
	c.JSON(http.StatusOK, common.OK)
}
