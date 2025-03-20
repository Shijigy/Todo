package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CheckinService 定义服务接口
type CheckinService interface {
	CreateCheckinService(ctx context.Context, checkin models.Checkin) (*models.CheckinWithDecodedCount, error)
	GetCheckinsByUserIDAndDateService(ctx context.Context, userID string, date string) ([]*models.CheckinWithDecodedCount, error)
	GetCheckinByIDService(ctx context.Context, checkinID int) (*models.CheckinWithDecodedCount, error)
	IncrementCheckinCountService(ctx context.Context, checkinID int, date string) (*models.Checkin, error)
	CheckIfCheckinCompleted(ctx context.Context, req models.CheckinCompletionRequest) (bool, error)
	UpdateCheckinService(ctx context.Context, checkinID int, checkin models.Checkin, startDate time.Time, endDate time.Time) (*models.CheckinWithDecodedCount, error)
	UpdateCountService(ctx context.Context, checkinID int, checkin models.Checkin, startDate string, endDate string, updatedCheckinCount []byte) error
	DeleteCheckinService(ctx context.Context, checkinID int) error
	ResetCheckinCountService(ctx context.Context, checkinID int, updatedCheckinCount []byte) error
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
func (s checkinService) CreateCheckinService(ctx context.Context, checkin models.Checkin) (*models.CheckinWithDecodedCount, error) {
	// 解析 start_date 和 end_date
	startDate, err := time.Parse("2006-01-02", checkin.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", checkin.EndDate)
	if err != nil {
		return nil, err
	}

	// 初始化打卡次数
	checkinCount := make(map[string]int)
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		checkinCount[dateStr] = 0 // 初始化每个日期的打卡次数为 0
	}

	// 将 checkinCount 序列化为 JSON 字符串
	checkinCountJSON, err := json.Marshal(checkinCount)
	if err != nil {
		return nil, err
	}

	// 设置 CheckinCount 为序列化后的 JSON 字符串
	checkin.CheckinCount = checkinCountJSON

	// 将数据保存到数据库
	err = s.repo.CreateCheckin(ctx, &checkin, checkin.CheckinCount) // 传递 JSON 格式的 checkin.CheckinCount
	if err != nil {
		return nil, err
	}

	// 解码 CheckinCount 字段的 JSON 字符串为 map[string]int
	var decodedCount map[string]int
	err = json.Unmarshal(checkin.CheckinCount, &decodedCount)
	if err != nil {
		return nil, err
	}

	// 创建返回的结构体，包含原始 Checkin 和解码后的数据
	result := &models.CheckinWithDecodedCount{
		Checkin:        &checkin,     // 返回原始的 Checkin 结构体
		DecodedCheckin: decodedCount, // 返回解码后的 CheckinCount 数据
	}

	return result, nil
}

// GetCheckinsByUserIDAndDateService 获取指定用户指定日期范围的打卡记录
func (s checkinService) GetCheckinsByUserIDAndDateService(ctx context.Context, userID string, date string) ([]*models.CheckinWithDecodedCount, error) {
	// 查询数据库，获取指定用户指定日期范围的打卡记录
	checkins, err := s.repo.GetCheckinsByUserIDAndDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkins by user_id and date: %v", err)
	}

	var result []*models.CheckinWithDecodedCount

	// 遍历每个打卡记录
	for _, checkin := range checkins {
		// 解码 CheckinCount 字段的 JSON 字符串为 map[string]int
		var decodedCount map[string]int
		err := json.Unmarshal(checkin.CheckinCount, &decodedCount)
		if err != nil {
			return nil, fmt.Errorf("failed to decode checkin_count for checkin ID %d: %v", checkin.ID, err)
		}

		// 创建返回的结构体，包含原始 Checkin 和解码后的数据
		result = append(result, &models.CheckinWithDecodedCount{
			Checkin:        &checkin,     // 返回原始的 Checkin 结构体
			DecodedCheckin: decodedCount, // 返回解码后的 CheckinCount 数据
		})
	}

	return result, nil
}

func (s checkinService) GetCheckinByIDService(ctx context.Context, checkinID int) (*models.CheckinWithDecodedCount, error) {
	checkins, err := s.repo.GetCheckinByID(ctx, checkinID)
	if err != nil {
		return nil, fmt.Errorf("获取打卡记录失败：%v", err)
	}

	// 解码 checkin_count
	var decodedCheckinCount map[string]int
	if checkins.CheckinCount != nil {
		err := json.Unmarshal(checkins.CheckinCount, &decodedCheckinCount)
		if err != nil {
			return nil, fmt.Errorf("解码 checkin_count 失败：%v", err)
		}
	}

	// 返回原始 Checkin 和解码后的 checkin_count
	return &models.CheckinWithDecodedCount{
		Checkin:        checkins,
		DecodedCheckin: decodedCheckinCount,
	}, nil
}

// CheckIfCheckinCompleted 判断是否完成打卡
func (s checkinService) CheckIfCheckinCompleted(ctx context.Context, req models.CheckinCompletionRequest) (bool, error) {
	// 获取打卡记录
	checkins, err := s.repo.GetCheckinByID(ctx, req.CheckinID)
	if err != nil {
		return false, fmt.Errorf("获取打卡记录失败：%v", err)
	}

	// 解码 checkin_count
	var decodedCheckinCount map[string]int
	if checkins.CheckinCount != nil {
		err := json.Unmarshal(checkins.CheckinCount, &decodedCheckinCount)
		if err != nil {
			return false, fmt.Errorf("解码 checkin_count 失败：%v", err)
		}
	}

	// 判断是否包含指定日期的打卡次数
	if count, exists := decodedCheckinCount[req.Date]; exists {
		// 如果当天的打卡次数等于目标次数，返回 true
		return count == checkins.TargetCheckinCount, nil
	}

	// 如果没有打卡记录，返回 false
	return false, nil
}

// IncrementCheckinCountService 增加指定日期的打卡次数
func (s *checkinService) IncrementCheckinCountService(ctx context.Context, checkinID int, date string) (*models.Checkin, error) {
	// 调用 repository 层的增量更新打卡次数
	checkin, err := s.repo.IncrementCheckinCount(ctx, checkinID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to increment checkin count: %w", err)
	}
	return checkin, nil
}

// UpdateCheckinService 更新打卡任务
func (s *checkinService) UpdateCheckinService(ctx context.Context, checkinID int, checkin models.Checkin, startDate time.Time, endDate time.Time) (*models.CheckinWithDecodedCount, error) {
	// 获取现有的打卡记录
	existingCheckin, err := s.repo.GetCheckinByID(ctx, checkinID)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkin: %v", err)
	}

	// 更新字段
	if checkin.Title != "" {
		existingCheckin.Title = checkin.Title
	}
	if checkin.StartDate != "" {
		existingCheckin.StartDate = checkin.StartDate
	}
	if checkin.EndDate != "" {
		existingCheckin.EndDate = checkin.EndDate
	}
	if checkin.Icon != 0 {
		existingCheckin.Icon = checkin.Icon
	}
	if checkin.TargetCheckinCount != 0 {
		existingCheckin.TargetCheckinCount = checkin.TargetCheckinCount
	}
	if checkin.MotivationMessage != "" {
		existingCheckin.MotivationMessage = checkin.MotivationMessage
	}

	// 重新计算 CheckinCount
	updatedCheckinCount := make(map[string]int)

	// 遍历更新后的日期范围
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		updatedCheckinCount[dateStr] = 0
	}

	// 保留已打卡的日期的打卡次数
	var decodedCount map[string]int
	err = json.Unmarshal(existingCheckin.CheckinCount, &decodedCount)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal existing checkin count: %v", err)
	}

	// 只保留新日期范围内的打卡记录
	for dateStr, count := range decodedCount {
		_, exists := updatedCheckinCount[dateStr]
		if exists {
			updatedCheckinCount[dateStr] = count
		}
	}

	// 序列化更新后的 checkinCount
	checkinCountJSON, err := json.Marshal(updatedCheckinCount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated checkin count: %v", err)
	}

	// 更新数据库中的 CheckinCount
	existingCheckin.CheckinCount = checkinCountJSON

	// 更新数据库中的记录
	err = s.repo.UpdateCheckin(ctx, existingCheckin)
	if err != nil {
		return nil, fmt.Errorf("failed to update checkin: %v", err)
	}

	// 解码 CheckinCount 字段的 JSON 字符串为 map[string]int
	var decodedUpdatedCount map[string]int
	err = json.Unmarshal(existingCheckin.CheckinCount, &decodedUpdatedCount)
	if err != nil {
		return nil, err
	}

	// 返回包含解码后的数据
	result := &models.CheckinWithDecodedCount{
		Checkin:        existingCheckin,     // 返回更新后的 Checkin 结构体
		DecodedCheckin: decodedUpdatedCount, // 返回更新后的打卡次数数据
	}

	return result, nil
}

// UpdateCountService 更新打卡记录
func (s checkinService) UpdateCountService(ctx context.Context, checkinID int, checkin models.Checkin, startDate string, endDate string, updatedCheckinCount []byte) error {
	// 更新打卡记录中的 checkin_count 字段
	err := s.repo.UpdateCount(ctx, checkinID, updatedCheckinCount)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCheckinService 删除打卡记录
func (s checkinService) DeleteCheckinService(ctx context.Context, checkinID int) error {
	return s.repo.DeleteCheckin(ctx, checkinID)
}

// ResetCheckinCountService 更新打卡次数
func (s checkinService) ResetCheckinCountService(ctx context.Context, checkinID int, updatedCheckinCount []byte) error {
	return s.repo.UpdateCount(ctx, checkinID, updatedCheckinCount)
}
