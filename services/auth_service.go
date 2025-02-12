package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"ToDo/utils"
	"context"
	"errors"
)

// UserService 用户服务接口
type UserService interface {
	RegisterUser(ctx context.Context, user models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

// 用户服务实现
type userService struct {
	repo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// RegisterUser 注册用户
func (s *userService) RegisterUser(ctx context.Context, user models.User) error {
	// 检查用户名是否已存在
	existingUser, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil && err != repositories.ErrUserNotFound {
		return err
	}
	if existingUser != nil {
		return errors.New("username already taken")
	}

	// 检查邮箱是否已存在
	existingEmailUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil && err != repositories.ErrUserNotFound {
		return err
	}
	if existingEmailUser != nil {
		return errors.New("email already registered")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// 保存用户到数据库
	_, err = s.repo.CreateUser(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// LoginUser 用户登录
func (s *userService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// 验证密码
	err = utils.ComparePasswordHash(password, user.Password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	// 生成 JWT Token
	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken 验证 Token 是否有效
func (s *userService) ValidateToken(ctx context.Context, token string) (string, error) {
	// 使用 utils 包的 ValidateJWTToken 方法来验证 Token
	username, err := utils.ValidateJWTToken(token)
	if err != nil {
		return "", err
	}

	return username, nil
}
