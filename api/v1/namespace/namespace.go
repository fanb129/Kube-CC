package namespace

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
	"net/http"
	"strconv"
)

// Index 展示所有namespace，
func Index(c *gin.Context) {
	g_id := c.DefaultQuery("g_id", "")
	u_id := c.DefaultQuery("u_id", "")

	var nsListResponse *responses.NsListResponse
	var err error
	// 1. u_id不为空，就是看指定用户的ns
	if u_id != "" {
		label := map[string]string{
			"u_id": u_id,
		}
		// 将map标签转换为string
		selector := labels.SelectorFromSet(label).String()
		nsListResponse, err = service.ListNs(selector)
		goto END
	}
	// 2. g_id不为空，u_id为空的话，就是查看该组下面所有人的ns
	if g_id != "" {
		var gid int
		gid, err = strconv.Atoi(g_id)
		if err != nil {
			zap.S().Errorln("ns:index:", err)
			goto END
		}
		users, err := dao.GetGroupUserById(uint(gid))

		nsListResponse = &responses.NsListResponse{
			Response: responses.OK,
		}
		for _, user := range users {
			label := map[string]string{
				"u_id": strconv.Itoa(int(user.ID)),
			}
			// 将map标签转换为string
			selector := labels.SelectorFromSet(label).String()
			var ns *responses.NsListResponse
			ns, err = service.ListNs(selector)
			if err != nil {
				zap.S().Errorln("ns:index:", err)
				goto END
			}
			// 拼接该组所有用户的ns
			nsListResponse.Length += ns.Length
			nsListResponse.NsList = append(nsListResponse.NsList, ns.NsList...)
		}
		goto END
	}
	//3. g_id和u_id都为空的话就是查看所有组下面所有人的ns
	nsListResponse, err = service.ListNs("")
	goto END

END:
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nsListResponse)
	}
	return
}

// Delete 删除指定namespace
func Delete(c *gin.Context) {
	ns := c.Param("ns")
	response, err := service.DeleteNs(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Add(c *gin.Context) {
	form := forms.NsAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	label := map[string]string{}
	if form.Uid != 0 {
		label["u_id"] = strconv.Itoa(int(form.Uid))
	}
	//expiredTime, err := time.Parse("2006-01-02 15:04:05", forms.ExpiredTime)
	//if err != nil {
	//	zap.S().Errorln(err)
	//	c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	//	return
	//}
	response, err := service.CreateNs(form.Name, form.ExpiredTime, label, form.Resources)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Update 更新namespace及其所含所有资源的uid
func Update(c *gin.Context) {
	form := forms.NsAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	uid := ""
	if form.Uid != 0 {
		uid = strconv.Itoa(int(form.Uid))
	}
	//expiredTime, err := time.Parse("2006-01-02 15:04:05", forms.ExpiredTime)
	//if err != nil {
	//	zap.S().Errorln(err)
	//	c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	//	return
	//}
	response, err := service.UpdateNs(form.Name, uid, form.ExpiredTime, form.Resources)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}
