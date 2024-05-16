package user

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"errors"
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
	g_id := c.Query("gid")
	gid, err := strconv.Atoi(g_id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	u_id, exists := c.Get("u_id")
	if !exists {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("获取权限信息失败").Error(),
		})
		return
	}
	uid := u_id.(uint)

	userListResponse, err := service.GetUserList(uint(gid), uid)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userListResponse)
}

//func GetAll(c *gin.Context) {
//	allUserList, err := service.GetAll()
//	if err != nil {
//		c.JSON(http.StatusOK, responses.Response{
//			StatusCode: -1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, allUserList)
//}

// // GetAd 获取管理员用户
// func GetAd(c *gin.Context) {
// 	response, err := service.GetAdminUser()
// 	if err != nil {
// 		c.JSON(http.StatusOK, responses.Response{
// 			StatusCode: -1,
// 			StatusMsg:  err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response)
// }

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

func Allocation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	form := forms.AllocationForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.AllocationUser(uint(id), form)
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

func SetEmail(c *gin.Context) {
	form := forms.SetEmailForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.SetEmail(form.Id, form.Email, form.VCode)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Add 添加注册用户
func Add(c *gin.Context) {
	form := forms.AddUserForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.RegisterUser(form.Username, form.Password, form.Nickname, form.Gid)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
