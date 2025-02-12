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

// UpdateCheckin 更新打卡记录
func (r *checkinRepository) UpdateCheckin(ctx context.Context, checkin *models.Checkin) (*models.Checkin, error) {
	if err := r.db.Save(checkin).Error; err != nil {
		return nil, err
	}
	return checkin, nil
}
