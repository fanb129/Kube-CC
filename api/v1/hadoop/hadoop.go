package hadoop

import (
	"Kube-CC/api/v1/namespace"
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/application"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

// Index 获取当前用户的Hadoop列表
//func Index(c *gin.Context) {
//	g_id := c.DefaultQuery("g_id", "")
//	u_id := c.DefaultQuery("u_id", "")
//
//	var hadoopListResponse *responses.HadoopListResponse
//	var err error
//	// 1. u_id不为空，就是看指定用户的ns
//	if u_id != "" {
//		hadoopListResponse, err = application.ListHadoop(u_id)
//		goto END
//	}
//	// 2. g_id不为空，u_id为空的话，就是查看该组下面所有人的ns
//	if g_id != "" {
//		var gid int
//		gid, err = strconv.Atoi(g_id)
//		if err != nil {
//			zap.S().Errorln("ns:index:", err)
//			goto END
//		}
//		users, err := dao.GetGroupUserById(uint(gid))
//
//		hadoopListResponse = &responses.HadoopListResponse{
//			Response: responses.OK,
//		}
//		for _, user := range users {
//			var hadoop *responses.HadoopListResponse
//			hadoop, err = application.ListHadoop(strconv.Itoa(int(user.ID)))
//			if err != nil {
//				zap.S().Errorln("ns:index:", err)
//				goto END
//			}
//			// 拼接该组所有用户的ns
//			hadoopListResponse.Length += hadoop.Length
//			hadoopListResponse.HadoopList = append(hadoop.HadoopList, hadoopListResponse.HadoopList...)
//		}
//		goto END
//	}
//	//3. g_id和u_id都为空的话就是查看所有组下面所有人的ns
//	hadoopListResponse, err = application.ListHadoop("")
//	goto END
//
//END:
//	if err != nil {
//		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
//	} else {
//		c.JSON(http.StatusOK, hadoopListResponse)
//	}
//	return
//}

// Index 根据查询条件列出hadoop列表
func Index(c *gin.Context) {
	namespace.BigDataIndex(c, application.ListHadoop)
}

// Add 创建hadoop
func Add(c *gin.Context) {
	// 表单验证
	form := forms.HadoopAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}

	res, err := application.CreateHadoop(form.Uid, form.Name, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas, form.ExpiredTime, form.ApplyResources)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
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
	res, err := application.DeleteHadoop(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func Update(c *gin.Context) {
	// 表单验证
	form := forms.HadoopUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	res, err := application.UpdateHadoop(form.Name, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas, form.ExpiredTime, form.ApplyResources)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
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
	form := forms.BatchHadoopAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	ids := form.Uid
	group := sync.WaitGroup{}
	group.Add(len(ids))
	for _, id := range ids {
		go func(id string) {
			if _, err := application.CreateHadoop(id, form.Name, form.HdfsMasterReplicas, form.DatanodeReplicas, form.YarnMasterReplicas, form.YarnNodeReplicas, form.ExpiredTime, form.ApplyResources); err != nil {
				zap.S().Errorln(err)
			}
			group.Done()
		}(id)
	}
	group.Wait()
	c.JSON(http.StatusOK, responses.OK)
}

func Info(c *gin.Context) {
	name := c.Query("name")
	response, err := application.GetHadoop(name)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
