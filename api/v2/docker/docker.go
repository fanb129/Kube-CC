package docker

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/docker"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 暂定先获取当前用户的全部镜像
func Info(c *gin.Context) {
	uid, ok := c.Get("u_id")
	if !ok {
		c.JSON(http.StatusOK, responses.NoSuchImage)
		return
	}

	rsp, err := docker.GetImage(uid.(string))

	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, rsp)
}

func Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	imageListResponse, err := docker.IndexDocker(page)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, imageListResponse)
	return
}

// Delete 删除镜像信息
func Remove(c *gin.Context) {
	id := c.Param("imageid")
	response, err := docker.DeleteImage(id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Update 更新tag
func Update(c *gin.Context) {
	id := c.Param("imageid")
	form := forms.ImageUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := docker.AddNewTag(id, form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

// TODO 后续完善备份镜像与拉取镜像的相关操作
// Save备份镜像
/*func Save(c *gin.Context) {
	id := c.Param("imageid")
	form := forms.SaveForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := docker.SaveImage(id)
}*/

// Pull拉取镜像
func Pull(c *gin.Context) {
	id := c.Param("imageid")
	form := forms.PullSpecifiedForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := docker.PullImage(id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
