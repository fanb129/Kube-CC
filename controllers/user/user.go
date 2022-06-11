package user

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/controllers"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/pkg/setting"
	"net/http"
	"strconv"
)

// Index 分页浏览用户信息
func Index(c *gin.Context) {
	//fmt.Println("userindex")
	page, _ := strconv.Atoi(c.Param("page"))
	total, userList := dao.GetUserList(page, setting.PageSize)

	// 如果无数据，则返回到第一页
	if total == 0 && page > 1 {
		page = 1
		total, userList = dao.GetUserList(page, setting.PageSize)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": total,
		"msg":  page,
		"data": userList,
	})

}

// Delete 删除用户
func Delete(c *gin.Context) {
	//fmt.Println("delete")
	id, _ := strconv.Atoi(c.Param("id"))
	user := dao.GetUserById(id)
	user.Status = 0
	dao.UpdateUser(user)
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "delete success",
		"data": user.ToMap(),
	})
}

// Edit 授权用户
func Edit(c *gin.Context) {
	//fmt.Println("useredit")
	id, _ := strconv.Atoi(c.Param("id"))
	newStatus, _ := strconv.Atoi(c.PostForm("status"))
	user := dao.GetUserById(id)
	user.Status = newStatus
	dao.UpdateUser(user)
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "edit success",
		"data": user.ToMap(),
	})
}

// ResetPass 重置密码
func ResetPass(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := dao.GetUserById(id)

	user.Password = controllers.EncryptionPWD(c.PostForm("password"))
	dao.UpdateUser(user)
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "reset success",
		"data": user.ToMap(),
	})
}
