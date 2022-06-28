package user

import (
	"github.com/gin-gonic/gin"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service"
	"net/http"
	"strconv"
)

func Index(c *gin.Context) {
	//fmt.Println("userindex")
	page, _ := strconv.Atoi(c.Param("page"))
	userListResponse := service.IndexUser(page)
	c.JSON(http.StatusOK, userListResponse)
}

// Delete 删除用户
func Delete(c *gin.Context) {
	//fmt.Println("delete")
	// 判断权限 管理员或者超级管理员
	jwt := GetStatus(c)
	if jwt >= 3 {
		id, _ := strconv.Atoi(c.Param("id"))
		response := service.DeleteUSer(id)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, common.NoStatus)
	}

}

// Edit 授权用户
func Edit(c *gin.Context) {
	//fmt.Println("useredit")
	// 判断权限 管理员或者超级管理员
	jwt := GetStatus(c)
	if jwt >= 3 {
		id, _ := strconv.Atoi(c.Param("id"))
		newStatus, _ := strconv.Atoi(c.PostForm("status"))
		response := service.EditUser(id, newStatus)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, common.NoStatus)
	}

}

// ResetPass 重置密码
func ResetPass(c *gin.Context) {
	jwt := GetStatus(c)
	if jwt >= 3 {
		id, _ := strconv.Atoi(c.Param("id"))
		password := c.PostForm("password")
		response := service.ResetPassUser(id, password)
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, common.NoStatus)
	}

}

// GetStatus 获得权限
func GetStatus(c *gin.Context) int {
	status, _ := c.Get("status")
	return status.(int)
}
