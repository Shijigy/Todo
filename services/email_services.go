package services

import (
	"net/smtp"
)

type EmailService interface {
	SendCaptcha(email, captchaCode string) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s *emailService) SendCaptcha(email, captchaCode string) error {
	from := "your-email@example.com"
	to := []string{email}
	subject := "Your Registration Captcha Code"
	body := "Your captcha code is: " + captchaCode

	msg := []byte("From: " + from + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	err := smtp.SendMail("smtp.example.com:587", nil, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
