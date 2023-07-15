package sys_login_log

import (
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPage 登录日志列表

func GetPage(c *gin.Context) {
	uid, ok := c.Get("u_id")
	if !ok {
		c.JSON(http.StatusOK, responses.NoUid)
		return
	}
	rsp, err := service.LogInfo(uid.(uint))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, rsp)
}

// Get 登录日志通过id获取

func Get(c *gin.Context) {

	//fmt.Println("userindex")
	page, _ := strconv.Atoi(c.Param("page"))
	userListResponse, err := service.IndexUser(page)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userListResponse)

}

// Delete 登录日志删除

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	response, err := service.DeleteUSer(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
