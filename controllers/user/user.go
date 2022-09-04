package user

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
	"strconv"
)

func Info(c *gin.Context) {
	uid, ok := c.Get("u_id")
	if !ok {
		c.JSON(http.StatusOK, common.NoUid)
		return
	}
	rsp, err := service.UserInfo(uid.(uint))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, rsp)
}
func Index(c *gin.Context) {
	//fmt.Println("userindex")
	page, _ := strconv.Atoi(c.Param("page"))
	userListResponse, err := service.IndexUser(page)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userListResponse)
}

// Delete 删除用户
func Delete(c *gin.Context) {
	//fmt.Println("delete")

	id, _ := strconv.Atoi(c.Param("id"))
	response, err := service.DeleteUSer(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)

}

// Edit 授权用户
func Edit(c *gin.Context) {
	//fmt.Println("useredit")
	id, _ := strconv.Atoi(c.Param("id"))
	form := common.EditForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	response, err := service.EditUser(uint(id), form.Role)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)

}

// ResetPass 重置密码
func ResetPass(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	form := common.ResetPassForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	response, err := service.ResetPassUser(uint(id), form.Password)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)

}
