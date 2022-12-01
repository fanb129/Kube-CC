package namespace

import (
	"Kube-CC/common"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"net/http"
	"strconv"
)

// Index 展示所有namespace，
func Index(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "")
	selector := ""
	if u_id != "" {
		label := map[string]string{
			"u_id": u_id,
		}
		// 将map标签转换为string
		selector = labels.SelectorFromSet(label).String()
	}

	nsListResponse, err := service.GetNs(selector)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nsListResponse)
	}
}

// Delete 删除指定namespace
func Delete(c *gin.Context) {
	ns := c.Param("ns")
	response, err := service.DeleteNs(ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Add(c *gin.Context) {
	form := common.NsAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	label := map[string]string{}
	if form.Uid != 0 {
		label["u_id"] = strconv.Itoa(int(form.Uid))
	}
	response, err := service.CreateNs(form.Name, label)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// Update 更新namespace及其所含所有资源的uid
func Update(c *gin.Context) {
	form := common.NsAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	uid := ""
	if form.Uid != 0 {
		uid = strconv.Itoa(int(form.Uid))
	}
	response, err := service.UpdateNs(form.Name, uid)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
