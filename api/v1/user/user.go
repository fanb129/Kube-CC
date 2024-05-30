package user

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"encoding/csv"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
	"strconv"
	"sync"

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

// AddByFile 通过文件批量添加
// csv 文件 username,nickname,password
func AddByFile(c *gin.Context) {
	form := forms.AddUserByFileForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	file, err := form.File.Open()
	//file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
	}
	if file != nil {
		defer file.Close()
	}
	//reader := csv.NewReader(file)
	//解决读取csv中文乱码的问题
	reader := csv.NewReader(transform.NewReader(file, simplifiedchinese.GBK.NewDecoder()))
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 忽略第一行的表头
	records = records[1:]

	// 协程批量添加
	var wg sync.WaitGroup
	for _, record := range records {
		wg.Add(1)
		go func(record []string) {
			defer wg.Done()
			// 从 record 中获取需要的数据
			username := record[0]
			nickname := record[1]
			password := record[2]
			// 使用 form.Gid 作为 group ID
			// 执行添加用户的逻辑
			if _, err = service.RegisterUser(username, password, nickname, form.Gid); err != nil {
				// 处理添加用户失败的情况
				zap.S().Errorln(err)
				return
			}
		}(record)
	}
	wg.Wait()
	c.JSON(http.StatusOK, &responses.OK)
}

func GetUserFile(c *gin.Context) {
	c.File("add_user.csv")
}
