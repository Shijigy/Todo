package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"
)

// UserService 用户服务接口
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	RegisterUser(ctx context.Context, user models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error) // 登录返回 token
}

// EmailService 邮箱服务接口
type EmailService interface {
	SendCaptcha(email, captchaCode string) error
}

// userService 用户服务实现
type userService struct {
	userRepo     repositories.UserRepository
	emailService EmailService
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository, emailService EmailService) UserService {
	return &userService{userRepo, emailService}
}

// GetUserByEmail 获取用户信息
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

// GenerateCaptcha 生成一个随机的验证码
func (s *userService) GenerateCaptcha() string {
	rand.Seed(time.Now().UnixNano())
	captcha := fmt.Sprintf("%06d", rand.Intn(1000000))
	return captcha
}

// RegisterUser 注册用户并发送验证码邮件
func (s *userService) RegisterUser(ctx context.Context, user models.User) error {
	// 1. 检查用户是否已经存在
	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// 2. 生成验证码
	captchaCode := s.GenerateCaptcha()

	// 3. 调用邮箱服务发送验证码
	err = s.emailService.SendCaptcha(user.Email, captchaCode)
	if err != nil {
		return fmt.Errorf("failed to send captcha: %v", err)
	}

	// 4. 保存用户到数据库
	_, err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}

	// 注册成功，无需进一步处理 createdUser
	return nil
}

// LoginUser 用户登录
func (s *userService) LoginUser(ctx context.Context, username, password string) (string, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(ctx, username)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// 验证密码
	// 此处假设密码验证已通过（通常会进行 bcrypt 验证）
	if password != user.Password {
		return "", fmt.Errorf("invalid credentials")
	}

	// 生成 token（你可以使用 JWT 等）
	token := "dummy-token" // 假设已经生成了一个 token
	return token, nil
}

// emailService 邮箱服务实现
type emailService struct {
	SMTPServer string
	FromEmail  string
	Password   string
}

// NewEmailService 创建邮箱服务实例
func NewEmailService(smtpServer, fromEmail, password string) EmailService {
	return &emailService{
		SMTPServer: smtpServer,
		FromEmail:  fromEmail,
		Password:   password,
	}
}

// SendCaptcha 发送验证码邮件
func (s *emailService) SendCaptcha(email, captchaCode string) error {
	from := s.FromEmail
	to := []string{email}
	subject := "Your Registration Captcha Code"
	body := "Your captcha code is: " + captchaCode

	// 构造邮件内容
	msg := []byte("From: " + from + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	// 设置SMTP服务器和身份验证
	auth := smtp.PlainAuth("", from, s.Password, s.SMTPServer)

	// 发送邮件
	err := smtp.SendMail(s.SMTPServer+":587", auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
