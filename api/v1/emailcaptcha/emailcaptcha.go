package emailcaptcha

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/dao"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetValidateCode(c *gin.Context) {
	// 获取目的邮箱
	em := []string{c.Param("email")}
	vCode, err := dao.SendEmailValidate(em)
	if err != nil {
		//log.Println(err)
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "发送验证码失败",
		})
		return
	}
	ctx := context.Background()
	// 验证码存入redis 并设置过期时间5分钟
	err = dao.RedisClient.Set(ctx, em[0], vCode, time.Duration(300*time.Second)).Err()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "验证码存储失败",
		})
		return
	}
	c.JSON(http.StatusOK, responses.Response{
		StatusCode: 1,
		StatusMsg:  "验证码存储成功",
	})
	return
}

func ValidateEmailCode(c *gin.Context) {
	//em := []string{c.Param("email")}
	// vCode := c.Param("vCode")
	verifyemailForm := forms.VerifyEmail{}
	// 参数绑定
	if err := c.ShouldBind(&verifyemailForm); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	em := verifyemailForm.Email
	vCode := verifyemailForm.VCode
	// 获取存储在redis中的验证码
	ctx := context.Background()
	vCodeRaw, err := dao.RedisClient.Get(ctx, em).Result()
	if err != nil {
		//log.Println(err.Error())
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "邮箱验证码已过期",
		})
		return
	}
	if vCodeRaw != "" && vCode == vCodeRaw {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: 1,
			StatusMsg:  "邮箱验证码验证成功",
		})
		return
	} else {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "邮箱验证码错误",
		})
		return
	}
}
