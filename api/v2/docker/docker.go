package docker

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/image"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func Index(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "0")
	g_id := c.DefaultQuery("g_id", "0")
	uid, err := strconv.Atoi(u_id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	gid, err := strconv.Atoi(g_id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	listResponse, err := image.ListImages(uint(uid), uint(gid))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, listResponse)
	}
}

func IndexOk(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "0")
	uid, err := strconv.Atoi(u_id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}

	listResponse, err := image.ListOkImages(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, listResponse)
	}
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	response, err := image.DeleteImage(uint(uid))
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Update(c *gin.Context) {
	form := forms.UpdateImageForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}

	response, err := image.UpdateImage(form)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Pull(c *gin.Context) {
	form := forms.PullImageForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}

	response, err := image.PullImage(form)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

//TODO
func Save(c *gin.Context) {

}
