package group

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//	func Info(c *gin.Context) {
//		gid, ok := c.Get("g_id")
//		if !ok {
//			c.JSON(http.StatusOK, responses.NoGid)
//			return
//		}
//		rsp, err := service.GroupInfo(gid.(uint))
//		if err != nil {
//			c.JSON(http.StatusOK, responses.Response{
//				StatusCode: -1,
//				StatusMsg:  err.Error(),
//			})
//		}
//		c.JSON(http.StatusOK, rsp)
//	}
func Index(c *gin.Context) {
	//fmt.Println("userindex")
	page, _ := strconv.Atoi(c.Param("page"))
	groupListResponse, err := service.IndexGroup(page)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, groupListResponse)
}

// Delete 删除组
func Delete(c *gin.Context) {
	//fmt.Println("delete")

	gid, _ := strconv.Atoi(c.Param("g_id"))
	response, err := service.DeleteGroup(uint(gid))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)

}

// ViewGroupUser 查看组内成员
func ViewGroupUser(c *gin.Context) {
	//fmt.Println("userindex")
	gid, _ := strconv.Atoi(c.Param("id"))
	groupuserListResponse, err := service.ViewGroupUser(uint(gid))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, groupuserListResponse)
}

// Remove 移出用户
func Remove(c *gin.Context) {
	//fmt.Println("useredit")
	uid, _ := strconv.Atoi(c.Param("u_id"))
	response, err := service.RemoveUser(uint(uid))
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
	uid, _ := strconv.Atoi(c.Param("u_id"))
	form := forms.GroupUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := service.UpdateGroup(uint(uid), form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)

}

// TransAdmin 更改管理员
func TransAdmin(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Param("g_id"))
	odid, _ := strconv.Atoi(c.Param("od_id"))
	nwid, _ := strconv.Atoi(c.Param("nw_id"))
	response, err := service.TransAdmin(uint(gid), uint(odid), uint(nwid))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
