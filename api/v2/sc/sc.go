package sc

import (
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index 展示所有sc
func Index(c *gin.Context) {
	scListRsp, err := service.ListSc()
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, scListRsp)
	}
}
