package utils

import (
	"fmt"
	"net/smtp"
)

// SendCaptcha 发送验证码邮件
func SendCaptcha(smtpServer, fromEmail, password, toEmail, captchaCode string) error {
	auth := smtp.PlainAuth("", fromEmail, password, smtpServer)
	to := []string{toEmail}
	subject := "Your Captcha Code"
	body := fmt.Sprintf("Your captcha code is: %s", captchaCode)

	// 构建邮件内容
	message := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// 发送邮件
	err := smtp.SendMail(smtpServer+":587", auth, fromEmail, to, message)
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}
	return nil
}
