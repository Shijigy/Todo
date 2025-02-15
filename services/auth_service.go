package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"ToDo/utils"
	"context"
	"fmt"
	"math/rand"
	"time"
)

// UserService 用户服务接口
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	RegisterUser(ctx context.Context, user models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error) // 登录返回 token
}

// userService 用户服务实现
type userService struct {
	userRepo     repositories.UserRepository
	emailService EmailService
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository, emailService EmailService) UserService {
	return &userService{
		userRepo:     userRepo,
		emailService: emailService,
	}
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
	// 检查用户是否已经存在
	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("用户邮箱 %s 已存在", user.Email)
	}

	// 生成验证码
	captchaCode := s.GenerateCaptcha()

	// 调用工具函数发送验证码
	err = utils.SendCaptcha(s.emailService.SMTPServer, s.emailService.FromEmail, s.emailService.Password, user.Email, captchaCode)
	if err != nil {
		return fmt.Errorf("发送验证码失败: %v", err)
	}

	// 保存用户到数据库
	_, err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return fmt.Errorf("用户注册失败: %v", err)
	}

	return nil
}

// LoginUser 用户登录
func (s *userService) LoginUser(ctx context.Context, username, password string) (string, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByEmail(ctx, username)
	if err != nil {
		return "", fmt.Errorf("凭证无效")
	}

	// 验证密码
	if password != user.Password {
		return "", fmt.Errorf("凭证无效")
	}

	// 生成 token（你可以使用 JWT 等）
	token := "dummy-token" // 假设已经生成了一个 token
	return token, nil
}
