package spark

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service/application"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

// Index 获取当前用户spark列表
func Index(c *gin.Context) {
	g_id := c.DefaultQuery("g_id", "")
	u_id := c.DefaultQuery("u_id", "")

	var sparkListResponse *responses.SparkListResponse
	var err error
	// 1. u_id不为空，就是看指定用户的ns
	if u_id != "" {
		sparkListResponse, err = application.ListSpark(u_id)
		goto END
	}
	// 2. g_id不为空，u_id为空的话，就是查看该组下面所有人的ns
	if g_id != "" {
		var gid int
		gid, err = strconv.Atoi(g_id)
		if err != nil {
			zap.S().Errorln("ns:index:", err)
			goto END
		}
		users, err := dao.GetGroupUserById(uint(gid))

		sparkListResponse = &responses.SparkListResponse{
			Response: responses.OK,
		}
		for _, user := range users {
			var spark *responses.SparkListResponse
			spark, err = application.ListSpark(strconv.Itoa(int(user.ID)))
			if err != nil {
				zap.S().Errorln("ns:index:", err)
				goto END
			}
			// 拼接该组所有用户的ns
			sparkListResponse.Length += spark.Length
			sparkListResponse.SparkList = append(spark.SparkList, spark.SparkList...)
		}
		goto END
	}
	//3. g_id和u_id都为空的话就是查看所有组下面所有人的ns
	sparkListResponse, err = application.ListSpark("")
	goto END

END:
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, sparkListResponse)
	}
	return
}

// Add 创建spark
func Add(c *gin.Context) {
	// 表单验证
	form := forms.SparkAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}

	res, err := application.CreateSpark(form.Uid, form.MasterReplicas, form.WorkerReplicas, form.ExpiredTime, form.ApplyResources)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func Delete(c *gin.Context) {
	ns := c.Param("name")
	response, err := application.DeleteSpark(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func Update(c *gin.Context) {
	// 表单验证
	form := forms.SparkUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	res, err := application.UpdateSpark(form.Name, form.MasterReplicas, form.WorkerReplicas, form.ExpiredTime, form.ApplyResources)
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
	form := forms.BatchSparkAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	ids := form.Uid
	group := sync.WaitGroup{}
	group.Add(len(ids))
	for _, id := range ids {
		go func(id string) {
			if _, err := application.CreateSpark(id, form.MasterReplicas, form.WorkerReplicas, form.ExpiredTime, form.ApplyResources); err != nil {
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
	response, err := application.GetSpark(name)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
