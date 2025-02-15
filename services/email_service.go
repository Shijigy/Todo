package services

// EmailService 邮件服务
type EmailService struct {
	SMTPServer string
	FromEmail  string
	Password   string
}

// NewEmailService 创建邮箱服务实例
func NewEmailService(smtpServer, fromEmail, password string) EmailService {
	return EmailService{
		SMTPServer: smtpServer,
		FromEmail:  fromEmail,
		Password:   password,
	}
}
