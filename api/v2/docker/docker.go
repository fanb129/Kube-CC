package docker

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/docker"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	uid, exists := c.Get("u_id")
	if exists == false {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "获取用户失败",
		})
		return
	}
	uId := uid.(uint)

	rsp, err := docker.IndexDocker(page, uId)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -100,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, rsp)
}

// Delete 删除镜像信息 测试通过
func Remove(c *gin.Context) {
	//id := c.Param("image_id")
	form := forms.RemoveImageForm{}

	//_, err, rsp := docker.DeleteImage(id)
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
	}
	_, err, rsp := docker.DeleteImage(form.ImageId)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -2,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rsp)
	//if err != nil {
	//	c.JSON(http.StatusOK, responses.Response{
	//		StatusCode: -1,
	//		StatusMsg:  err.Error(),
	//	})
	//} else {
	//	c.JSON(http.StatusOK, rsp)
	//}
}

// ADD 通过新增tag对镜像进行创建 测试通过
func TagAdd(c *gin.Context) {
	form := forms.ImageCreateByTagForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := docker.AlertTag(form.OldRepositoryName, form.OldTag, form.NewRepositoryName, form.NewTag)
	if err != nil {
		zap.S().Errorln(err)
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

// Pull拉取镜像 测试通过
func PullPublic(c *gin.Context) {
	form := forms.PullFromRepositoryPublicForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := docker.PullImage(form.RepositoryName, form.Tag, form.Uid, form.Kind)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rsp)
}

// 拉取私有仓库的镜像，除新增账号密码参数以外与pullpublic一致 预计通过
func PullPrivate(c *gin.Context) {
	form := forms.PullFromRepositoryPrivateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := docker.PullPrivateImage(form.RepositoryName, form.Tag, form.Username, form.Passwd, form.Uid, form.Kind)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func CreateImageByImageId(c *gin.Context) {
	form := forms.ImageCreateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := docker.CreateImage(form.Parent, form.Username, form.Passwd, form.Tag, form.Uid, form.Kind)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rsp)
}
