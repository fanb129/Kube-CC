package dao

import (
	"fmt"
	"math/rand" //随机数
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

func SendEmailValidate(em string) (string, error) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Kube-CC <1916861581@qq.com>")
	e.To = []string{em}
	// 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05")
	//设置文件发送的内容
	content := fmt.Sprintf(`
		您本次的验证码为%s
		尊敬的%s，您好！您于 %s 提交本次邮箱验证，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。
		此邮箱为系统邮箱，请勿回复。
	`, vCode, em, t)
	e.Text = []byte(content)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "1916861581@qq.com", "rxtspyuerwocbcae", "smtp.qq.com"))
	return vCode, err
}
