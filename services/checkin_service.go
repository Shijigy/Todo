package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"errors"
)

// CheckinService 定义服务接口
type CheckinService interface {
	CreateCheckinService(ctx context.Context, checkin models.Checkin) (*models.Checkin, error)
	GetCheckinRecords(ctx context.Context, userID string) ([]models.Checkin, error)
}

// 服务实现
type checkinService struct {
	repo repositories.CheckinRepository
}

// NewCheckinService 创建服务实例
func NewCheckinService(repo repositories.CheckinRepository) CheckinService {
	return &checkinService{repo: repo}
}

// CreateCheckinService 创建打卡任务
func (s *checkinService) CreateCheckinService(ctx context.Context, checkin models.Checkin) (*models.Checkin, error) {
	createdCheckin, err := s.repo.CreateCheckin(ctx, &checkin)
	if err != nil {
		return nil, err
	}
	return createdCheckin, nil
}

// GetCheckinRecords 获取用户的打卡记录
func (s *checkinService) GetCheckinRecords(ctx context.Context, userID string) ([]models.Checkin, error) {
	checkins, err := s.repo.GetCheckinByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("no check-in records found")
	}
	return checkins, nil
}
