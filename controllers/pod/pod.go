package pod

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
)

func GetPod(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "")
	ns := c.DefaultQuery("namespace", "")
	selector := ""
	if u_id != "" {
		label := map[string]string{
			"u_id": u_id,
		}
		// 将map标签转换为string
		selector = labels.SelectorFromSet(label).String()
	}

	fmt.Println(u_id, selector, ns)
	podListResponse, err := service.GetPod(ns, selector)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, podListResponse)
	}
}
