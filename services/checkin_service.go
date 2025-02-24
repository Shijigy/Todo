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
	GetCheckinRecordsByUserID(ctx context.Context, userID string) ([]models.Checkin, error) // 修改返回值类型为 []models.Checkin
	MarkCheckinCompleteService(ctx context.Context, checkinID string) (*models.Checkin, error)
	UpdateCheckinCountService(ctx context.Context, checkinID string, increment int) (*models.Checkin, error)
	DeleteCheckinService(ctx context.Context, checkinID string) error
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

// GetCheckinRecordsByUserID 获取用户所有的打卡记录（去掉taskID）
func (s *checkinService) GetCheckinRecordsByUserID(ctx context.Context, userID string) ([]models.Checkin, error) {
	checkins, err := s.repo.GetCheckinByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("no check-in records found for the given user")
	}
	return checkins, nil
}

// MarkCheckinCompleteService 标记打卡任务完成
func (s *checkinService) MarkCheckinCompleteService(ctx context.Context, checkinID string) (*models.Checkin, error) {
	checkin, err := s.repo.GetCheckinByID(ctx, checkinID)
	if err != nil {
		return nil, err
	}

	// 增加打卡次数
	checkin.CheckinCount++

	// 如果打卡次数达到目标次数，更新状态为 completed
	if checkin.CheckinCount >= checkin.TargetCheckinCount {
		checkin.Status = "completed"
	}

	updatedCheckin, err := s.repo.UpdateCheckin(ctx, checkin)
	if err != nil {
		return nil, err
	}
	return updatedCheckin, nil
}

// UpdateCheckinCountService 更新打卡次数
func (s *checkinService) UpdateCheckinCountService(ctx context.Context, checkinID string, increment int) (*models.Checkin, error) {
	// 获取当前打卡记录
	checkin, err := s.repo.GetCheckinByID(ctx, checkinID)
	if err != nil {
		return nil, errors.New("checkin not found")
	}

	// 增加打卡次数
	checkin.CheckinCount += increment

	// 如果打卡次数达到了目标次数，更新状态为 "completed"
	if checkin.CheckinCount >= checkin.TargetCheckinCount {
		checkin.Status = "completed"
	}

	// 更新数据库中的打卡记录
	updatedCheckin, err := s.repo.UpdateCheckin(ctx, checkin)
	if err != nil {
		return nil, err
	}

	return updatedCheckin, nil
}

// 删除打卡任务
func (s *checkinService) DeleteCheckinService(ctx context.Context, checkinID string) error {
	// 调用仓库层删除打卡记录
	err := s.repo.DeleteCheckin(ctx, checkinID)
	if err != nil {
		return err
	}
	return nil
}
