package login

import (
	"Kube-CC/common/forms"
	"Kube-CC/common/responses"
	"Kube-CC/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// CheckPass 验证密码
func CheckPass(c *gin.Context) {
	loginForm := forms.LoginForm{}
	// 参数绑定
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}

	// 调用业务层登录
	loginRes, err := service.Login(loginForm.UsernameorEmail, loginForm.Password)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: 0, // 返回0，前端不弹出错误提示框
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginRes)
}

// Login 用户登录
func Login(c *gin.Context) {
	loginForm := forms.LoginForm{}
	// 参数绑定
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}

	// 调用业务层登录
	loginRes, err := service.Login(loginForm.UsernameorEmail, loginForm.Password)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginRes)
}

func Forget(c *gin.Context) {
	c.JSON(http.StatusOK, responses.Response{
		StatusCode: 1,
		StatusMsg:  responses.OK.StatusMsg,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	c.JSON(200, responses.OK)
}

// Register 用户注册
func Register(c *gin.Context) {
	//fmt.Println("register")
	// 表单验证，参数绑定
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	registerRes, err := service.Register(registerForm.Username, registerForm.Password, registerForm.Nickname, registerForm.Email)
	if err != nil {
		c.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, registerRes)
}

var store = base64Captcha.DefaultMemStore

// GetCaptcha 生成验证码
func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(60, 200, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		ctx.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "生成验证码失败",
		})
		zap.S().Errorln("生成验证码错误: ", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, responses.CaptchaResponse{
		Response:  responses.OK,
		CaptchaId: id,
		PicPath:   b64s,
	})
}

// CheckCaptcha 验证验证码
func CheckCaptcha(ctx *gin.Context) {
	verifycaptchaForm := forms.VerifyCaptcha{}
	if err := ctx.ShouldBind(&verifycaptchaForm); err != nil {
		ctx.JSON(http.StatusOK, responses.ValidatorResponse(err))
		return
	}
	if verifycaptchaForm.CaptchaVal == "" {
		ctx.JSON(http.StatusOK, responses.Response{
			StatusCode: 0,
			StatusMsg:  "未输入图片验证码",
		})
		return
	}
	res := store.Verify(verifycaptchaForm.CaptchaId, verifycaptchaForm.CaptchaVal, true)
	if !res {
		ctx.JSON(http.StatusOK, responses.Response{
			StatusCode: -1,
			StatusMsg:  "图片验证码错误",
		})
		return
	}
	// id := ctx.DefaultQuery("captchaId", "")
	// captchaVal := ctx.DefaultQuery("captchaVal", "")
	// res := store.Verify(id, captchaVal, true)
	ctx.JSON(http.StatusOK, responses.Response{
		StatusCode: 1,
		StatusMsg:  "图片验证码正确",
	})
}
