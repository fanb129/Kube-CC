package captcha

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var (
	// 存储方式
	store = base64Captcha.DefaultMemStore
	// 图片高度
	height = 40
	// 图片宽度
	width = 240
	// 干扰线选项
	showLineOptions = 0 // 2 4 8
	// 干扰数
	noiseCount = 0
	// 验证码长度
	length = 4
	// 字体
	fonts = []string{"wqy-microhei.ttc"}
	// 字符
	stringSource = "1234567890QAZWSXCDERFVBGTYHNMJUIKLOP"
	// 中文
	chineseSource = "你,好,真,香,哈,哈,消,费,者,狗,仔,北京烤鸭"
	// 音频语言
	language = "zh"
	// 最大绝对偏斜系数为个位数
	maxSkew = 0
	// 圆圈的数量
	dotCount = 1
)

// GetCaptcha 获取验证码
func GetCaptcha(ctx *gin.Context) {
	id, b64s, err := CaptchaDigit()
	if err != nil {
		zap.S().Errorln("验证码获取失败:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "验证码获取失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success", "data": map[string]interface{}{
		"id":      id,
		"captcha": b64s,
	}})
	return
}

// CheckCaptcha 验证验证码
func CheckCaptcha(ctx *gin.Context) {
	id := ctx.DefaultQuery("id", "")
	value := ctx.DefaultQuery("value", "")
	if id == "" || "" == value {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数不完整"})
		return
	}
	b := CaptchaVerify(id, value)
	if !b {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "验证码错误"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
	return
}

// CaptchaString 生成字符串验证码
func CaptchaString() (id, b64s string, err error) {
	// 字符串验证码
	driver := &base64Captcha.DriverString{
		Height:          height,
		Width:           width,
		ShowLineOptions: showLineOptions,
		NoiseCount:      noiseCount,
		Source:          stringSource,
		Length:          length,
		Fonts:           fonts,
	}
	driver = driver.ConvertFonts()
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

// CaptchaMath 生成算术验证码
func CaptchaMath() (id, b64s string, err error) {
	driver := &base64Captcha.DriverMath{
		Height:          height,
		Width:           width,
		NoiseCount:      noiseCount,
		ShowLineOptions: showLineOptions,
		Fonts:           fonts,
	}
	driver = driver.ConvertFonts()
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

// CaptchaChinese 生成中文验证码
func CaptchaChinese() (id, b64s string, err error) {
	driver := &base64Captcha.DriverChinese{
		Height:          height,
		Width:           width,
		NoiseCount:      noiseCount,
		ShowLineOptions: showLineOptions,
		Length:          length,
		Source:          chineseSource,
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	driver = driver.ConvertFonts()
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

// CaptchaAudio 生成音频验证码
func CaptchaAudio() (id, b64s string, err error) {
	driver := &base64Captcha.DriverAudio{
		Length: length,
		// "en", "ja", "ru", "zh".
		Language: language,
	}
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

// CaptchaDigit 生成数字验证码
func CaptchaDigit() (id, b64s string, err error) {
	driver := &base64Captcha.DriverDigit{
		Height: height,
		Width:  width,
		Length: length,
		// 最大绝对偏斜系数为个位数
		MaxSkew: 0,
		// 背景圆圈的数量。
		DotCount: 0,
	}
	return base64Captcha.NewCaptcha(driver, store).Generate()
}

func CaptchaVerify(id, value string) bool {
	return store.Verify(id, value, true)
}
