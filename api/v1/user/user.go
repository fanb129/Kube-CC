package user

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Info(c *gin.Context) {
	uid, ok := c.Get("u_id")
	if !ok {
		c.JSON(http.StatusOK, responses.NoUid)
		return
	}
	rsp, err := service.UserInfo(uid.(uint))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
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
		c.JSON(http.StatusOK, responses.Response{
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
		c.JSON(http.StatusOK, responses.Response{
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
	form := forms.EditForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.EditUser(uint(id), form.Role)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)

}

// Update 更新用户信息
func Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	form := forms.UpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.UpdateUser(uint(id), form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
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
	form := forms.ResetPassForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.ResetPassUser(uint(id), form.Password)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
