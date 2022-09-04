package hadoop

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
)

// Index 获取当前用户的Hadoop列表
func Index(c *gin.Context) {
	u_id, ok := c.Get("u_id")
	if !ok {
		c.JSON(http.StatusOK, common.NoUid)
		return
	}
	hadoopListResponse, err := service.GetHadoop(u_id.(uint))
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

	res, err := service.CreateHadoop(form.Uid, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas)
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
