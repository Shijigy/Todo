package repositories

import (
	"ToDo/models"
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// CheckinRepository 打卡任务仓库接口
type CheckinRepository interface {
	CreateCheckin(ctx context.Context, checkin *models.Checkin) (*models.Checkin, error)
	GetCheckinByUserID(ctx context.Context, userID string) ([]models.Checkin, error)
	UpdateCheckin(ctx context.Context, checkin *models.Checkin) (*models.Checkin, error)
	GetCheckinByID(ctx context.Context, checkinID string) (*models.Checkin, error)
	DeleteCheckin(ctx context.Context, checkinID string) error
}

// 打卡任务仓库实现
type checkinRepository struct {
	db *gorm.DB
}

// NewCheckinRepository 创建打卡任务仓库实例
func NewCheckinRepository(db *gorm.DB) CheckinRepository {
	return &checkinRepository{db: db}
}

// CreateCheckin 创建打卡记录
func (r *checkinRepository) CreateCheckin(ctx context.Context, checkin *models.Checkin) (*models.Checkin, error) {
	// 获取用户当前的打卡次数
	var count int
	err := r.db.Model(&models.Checkin{}).Where("user_id = ?", checkin.UserID).Count(&count).Error
	if err != nil {
		return nil, err
	}

	// 更新打卡次数
	checkin.CheckinCount = count + 1

	// 创建新的打卡记录
	checkin.CheckinAt = time.Now()
	if err := r.db.Create(checkin).Error; err != nil {
		return nil, err
	}
	return checkin, nil
}

// GetCheckinByUserID 根据用户 ID 获取打卡记录
func (r *checkinRepository) GetCheckinByUserID(ctx context.Context, userID string) ([]models.Checkin, error) {
	var checkins []models.Checkin
	if err := r.db.Where("user_id = ?", userID).Find(&checkins).Error; err != nil {
		return nil, err
	}
	return checkins, nil
}

// GetCheckinByTaskIDAndUserID 获取用户某个任务的打卡记录
func (r *checkinRepository) GetCheckinByTaskIDAndUserID(ctx context.Context, userID, taskID string) (*models.Checkin, error) {
	var checkin models.Checkin
	if err := r.db.Where("user_id = ? AND task_id = ?", userID, taskID).First(&checkin).Error; err != nil {
		return nil, err
	}
	return &checkin, nil
}

// GetCheckinByID 根据 checkinID 获取打卡记录
func (r *checkinRepository) GetCheckinByID(ctx context.Context, checkinID string) (*models.Checkin, error) {
	var checkin models.Checkin
	if err := r.db.Where("id = ?", checkinID).First(&checkin).Error; err != nil {
		return nil, err
	}
	return &checkin, nil
}

// UpdateCheckin 更新打卡记录
func (r *checkinRepository) UpdateCheckin(ctx context.Context, checkin *models.Checkin) (*models.Checkin, error) {
	if err := r.db.Save(checkin).Error; err != nil {
		return nil, err
	}
	return checkin, nil
}

// DeleteCheckin 删除打卡任务
func (r *checkinRepository) DeleteCheckin(ctx context.Context, checkinID string) error {
	// 执行删除操作
	if err := r.db.Where("id = ?", checkinID).Delete(&models.Checkin{}).Error; err != nil {
		return err
	}
	return nil
}
