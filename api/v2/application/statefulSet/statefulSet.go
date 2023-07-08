package statefulSet

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service/application"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Add(c *gin.Context) {
	form := forms.StatefulSetAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	rsp, err := application.CreateAppStatefulSet(form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.DeleteAppSetfulset(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	if ns == "" {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: "请选择namespace"})
	}
	response, err := application.ListAppStatesulSet(ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Update(c *gin.Context) {
	form := forms.StatefulSetAddForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	response, err := application.UpdateAppStatefulSet(form)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Info(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := application.GetAppStatefulSet(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

//
//// Index 展示所有statefulSet
//func Index(c *gin.Context) {
//	ns := c.DefaultQuery("ns", "")
//	/*u_id := c.DefaultQuery("u_id", "")
//	selector := ""
//	if u_id != "" {
//		label := map[string]string{
//			"u_id": u_id,
//		}
//		// 将map标签转换为string
//		selector = labels.SelectorFromSet(label).String()
//	}*/
//	statefulSetListResponse, err := service.GetStatefulSet(ns, "")
//	if err != nil {
//		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
//	} else {
//		c.JSON(http.StatusOK, statefulSetListResponse)
//	}
//}
//
//// Delete 删除指定statefulSet
//func Delete(c *gin.Context) {
//	ns := c.Query("ns")
//	name := c.Query("name")
//	response, err := service.DeleteStatefulSet(name, ns)
//	if err != nil {
//		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
//	} else {
//		c.JSON(http.StatusOK, response)
//	}
//}
//
//// Info 获取单个statefulSet的yaml信息
//func Info(c *gin.Context) {
//	ns := c.Query("ns")
//	name := c.Query("name")
//	res, err := service.GetAStatefulSet(name, ns)
//	if err != nil {
//		c.JSON(http.StatusOK, responses.Response{StatusCode: -1, StatusMsg: err.Error()})
//	} else {
//		c.JSON(http.StatusOK, responses.StatefulSetInfo{
//			Response: responses.OK,
//			Info:     *res,
//		})
//	}
//}
//
//// Add 通过表单提交添加StatefulSet
////func Add(c *gin.Context) {
////	v1.statefulSet{}
////	yamlApply.DeployCreate()
////}
