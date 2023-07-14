package namespace

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"Kube-CC/service"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/uuid"
	"net/http"
	"strconv"
)

// Index 展示所有的工作空间 namespace，
func Index(c *gin.Context) {
	label := map[string]string{
		"kind": "workspace",
	}
	g_id := c.DefaultQuery("g_id", "")
	u_id := c.DefaultQuery("u_id", "")

	// 1. u_id不为空，就是看指定用户的ns
	if u_id != "" {
		intuid, err := strconv.Atoi(u_id)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		uintuid := uint(intuid)
		// 判断权限
		// 是否为自己的ns
		uid, exists := c.Get("u_id")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取权限信息失败").Error(),
			})
			return
		}
		if uid.(uint) == uintuid {
			label["u_id"] = u_id
			// 将map标签转换为string
			selector := labels.SelectorFromSet(label).String()
			nsListResponse, err := service.ListNs(selector)
			if err != nil {
				c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			} else {
				c.JSON(http.StatusOK, nsListResponse)
			}
			return
		}
		// 是否为超管
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取权限信息失败").Error(),
			})
			return
		}
		if role.(uint) == 3 {
			label["u_id"] = u_id
			// 将map标签转换为string
			selector := labels.SelectorFromSet(label).String()
			nsListResponse, err := service.ListNs(selector)
			if err != nil {
				c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			} else {
				c.JSON(http.StatusOK, nsListResponse)
			}
			return
		}
		// 是否为自己组员
		user, err := dao.GetUserById(uintuid)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		group, err := dao.GetGroupById(user.Groupid)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		if group.Adminid == uid.(uint) {
			label["u_id"] = u_id
			// 将map标签转换为string
			selector := labels.SelectorFromSet(label).String()
			nsListResponse, err := service.ListNs(selector)
			if err != nil {
				c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			} else {
				c.JSON(http.StatusOK, nsListResponse)
			}
			return
		}

		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("没有权限").Error(),
		})
		return
	}
	// 2. g_id不为空，u_id为空的话，就是查看该组下面所有人的ns
	if g_id != "" {
		// 判断权限
		uid, exists := c.Get("u_id")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取用户信息失败").Error(),
			})
			return
		}
		gid, err := strconv.Atoi(g_id)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}
		group, err := dao.GetGroupById(uint(gid))
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取权限信息失败").Error(),
			})
			return
		}
		if group.Adminid != uid.(uint) && role.(uint) != 3 {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("没有权限").Error(),
			})
			return
		}
		users, err := dao.GetGroupUserById(uint(gid))
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}

		nsListResponse := &responses.NsListResponse{
			Response: responses.OK,
		}
		for _, user := range users {
			label["u_id"] = strconv.Itoa(int(user.ID))
			// 将map标签转换为string
			selector := labels.SelectorFromSet(label).String()
			var ns *responses.NsListResponse
			ns, err = service.ListNs(selector)
			if err != nil {
				zap.S().Errorln("ns:index:", err)
				c.JSON(http.StatusOK, responses.Response{
					StatusCode: -1,
					StatusMsg:  err.Error(),
				})
				return
			}
			// 拼接该组所有用户的ns
			nsListResponse.Length += ns.Length
			nsListResponse.NsList = append(nsListResponse.NsList, ns.NsList...)
		}
		c.JSON(http.StatusOK, nsListResponse)
		return
	}
	//3. g_id和u_id都为空的话就是查看所有组下面所有人的ns
	// 判断权限
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("获取权限信息失败").Error(),
		})
		return
	}
	if role.(uint) != 3 {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("没有权限").Error(),
		})
		return
	}
	// 将map标签转换为string
	selector := labels.SelectorFromSet(label).String()
	nsListResponse, err := service.ListNs(selector)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nsListResponse)
	}
}

// Delete 删除指定namespace
func Delete(c *gin.Context) {
	ns := c.Param("ns")
	get, err := service.GetNs(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	nsUid := get.Labels["u_id"]
	intuid, err := strconv.Atoi(nsUid)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	uintuid := uint(intuid)
	// 判断权限
	// 是否为自己的ns
	uid, exists := c.Get("u_id")
	if !exists {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("获取权限信息失败").Error(),
		})
		return
	}
	if uid.(uint) == uintuid {
		response, err := service.DeleteNs(ns)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		} else {
			c.JSON(http.StatusOK, response)
		}
	}
	// 是否为超管
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("获取权限信息失败").Error(),
		})
		return
	}
	if role.(uint) == 3 {
		response, err := service.DeleteNs(ns)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		} else {
			c.JSON(http.StatusOK, response)
		}
	}
	// 是否为自己组员的ns
	user, err := dao.GetUserById(uintuid)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	group, err := dao.GetGroupById(user.Groupid)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	if group.Adminid == uid.(uint) {
		response, err := service.DeleteNs(ns)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		} else {
			c.JSON(http.StatusOK, response)
		}
		return
	}
	c.JSON(http.StatusOK, responses.Response{
		StatusCode: -1,
		StatusMsg:  errors.New("没有权限").Error(),
	})
}

// Add 创建workspace
func Add(c *gin.Context) {
	form := forms.NsAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	label := map[string]string{
		"kind": "workspace",
	}
	if form.Uid != 0 {
		label["u_id"] = strconv.Itoa(int(form.Uid))
	}
	//expiredTime, err := time.Parse("2006-01-02 15:04:05", forms.ExpiredTime)
	//if err != nil {
	//	zap.S().Errorln(err)
	//	c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	//	return
	//}
	newUUID := string(uuid.NewUUID())
	response, err := service.CreateNs(form.Name+"-"+newUUID, "", label, form.Resources)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Update 更新namespace
func Update(c *gin.Context) {
	form := forms.NsUpdateForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	//expiredTime, err := time.Parse("2006-01-02 15:04:05", forms.ExpiredTime)
	//if err != nil {
	//	zap.S().Errorln(err)
	//	c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	//	return
	//}
	response, err := service.UpdateNs(form.Name, "", form.Resources)
	if err != nil {
		zap.S().Errorln(err)
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// BigDataIndex  将spark，hadoop的index操作封装在一起
func BigDataIndex(c *gin.Context, listFun func(uid string) (*responses.BigdataListResponse, error)) {
	g_id := c.DefaultQuery("g_id", "")
	u_id := c.DefaultQuery("u_id", "")

	// 1. u_id不为空，就是看指定用户的ns
	if u_id != "" {
		intuid, err := strconv.Atoi(u_id)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		uintuid := uint(intuid)
		// 判断权限
		// 是否为自己的ns
		uid, exists := c.Get("u_id")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取权限信息失败").Error(),
			})
			return
		}
		if uid.(uint) == uintuid {
			nsListResponse, err := listFun(u_id)
			if err != nil {
				c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			} else {
				c.JSON(http.StatusOK, nsListResponse)
			}
			return
		}
		// 是否为超管
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取权限信息失败").Error(),
			})
			return
		}
		if role.(uint) == 3 {
			nsListResponse, err := listFun(u_id)
			if err != nil {
				c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			} else {
				c.JSON(http.StatusOK, nsListResponse)
			}
			return
		}
		// 是否为自己组员
		user, err := dao.GetUserById(uintuid)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		group, err := dao.GetGroupById(user.Groupid)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			return
		}
		if group.Adminid == uid.(uint) {
			nsListResponse, err := listFun(u_id)
			if err != nil {
				c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
			} else {
				c.JSON(http.StatusOK, nsListResponse)
			}
			return
		}

		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("没有权限").Error(),
		})
		return
	}
	// 2. g_id不为空，u_id为空的话，就是查看该组下面所有人的ns
	if g_id != "" {
		// 判断权限
		uid, exists := c.Get("u_id")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取用户信息失败").Error(),
			})
			return
		}
		gid, err := strconv.Atoi(g_id)
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}
		group, err := dao.GetGroupById(uint(gid))
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("获取权限信息失败").Error(),
			})
			return
		}
		if group.Adminid != uid.(uint) && role.(uint) != 3 {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  errors.New("没有权限").Error(),
			})
			return
		}
		users, err := dao.GetGroupUserById(uint(gid))
		if err != nil {
			c.JSON(http.StatusOK, responses.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			return
		}

		nsListResponse := &responses.BigdataListResponse{
			Response: responses.OK,
		}
		for _, user := range users {
			var ns *responses.BigdataListResponse
			ns, err = listFun(strconv.Itoa(int(user.ID)))
			if err != nil {
				zap.S().Errorln("ns:index:", err)
				c.JSON(http.StatusOK, responses.Response{
					StatusCode: -1,
					StatusMsg:  err.Error(),
				})
				return
			}
			// 拼接该组所有用户的ns
			nsListResponse.Length += ns.Length
			nsListResponse.BigdataList = append(nsListResponse.BigdataList, ns.BigdataList...)
		}
		c.JSON(http.StatusOK, nsListResponse)
		return
	}
	//3. g_id和u_id都为空的话就是查看所有组下面所有人的ns
	// 判断权限
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("获取权限信息失败").Error(),
		})
		return
	}
	if role.(uint) != 3 {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  errors.New("没有权限").Error(),
		})
		return
	}

	nsListResponse, err := listFun("")
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, nsListResponse)
	}
}

func NsTotal(c *gin.Context) {
	u_id := c.DefaultQuery("u_id", "")
	if u_id == "" {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: "获取uid失败"})
		return
	}
	rsp, err := service.GetUserNsTotal(u_id)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, rsp)
}
