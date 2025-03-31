package repositories

import (
	"ToDo/models"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// CheckinRepository 打卡任务仓库接口
type CheckinRepository interface {
	CreateCheckin(ctx context.Context, checkin *models.Checkin, checkinCountJSON []byte) error
	GetCheckinsByUserIDAndDate(ctx context.Context, userID string, date string) ([]models.Checkin, error)
	GetCheckinByID(ctx context.Context, checkinID int) (*models.Checkin, error)
	IncrementCheckinCount(ctx context.Context, checkinID int, date string) (*models.CheckinWithDecodedCount, error)
	UpdateCheckin(ctx context.Context, checkin *models.Checkin) error
	UpdateCount(ctx context.Context, checkinID int, updatedCheckinCount []byte) error
	DeleteCheckin(ctx context.Context, checkinID int) error
}

// 打卡任务仓库实现
type checkinRepository struct {
	db *gorm.DB
}

// NewCheckinRepository 创建打卡任务仓库实例
func NewCheckinRepository(db *gorm.DB) CheckinRepository {
	return &checkinRepository{
		db: db,
	}
}

// CreateCheckin 创建打卡记录到数据库
func (r *checkinRepository) CreateCheckin(ctx context.Context, checkin *models.Checkin, checkinCountJSON []byte) error {
	// 插入数据库时将 checkinCountJSON 作为一个 JSON 字符串存入数据库
	query := `
		INSERT INTO checkins (
			user_id, 
			title, 
			start_date, 
			end_date, 
			checkin_count, 
			target_checkin_count, 
			icon, 
			motivation_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// 执行插入操作
	if err := r.db.Exec(query,
		checkin.UserID,
		checkin.Title,
		checkin.StartDate,
		checkin.EndDate,
		checkinCountJSON, // 这里直接传入 checkinCountJSON
		checkin.TargetCheckinCount,
		checkin.Icon,
		checkin.MotivationMessage,
	).Error; err != nil {
		return err
	}

	return nil
}

// GetCheckinsByUserIDAndDate 获取指定用户指定日期的打卡记录
func (r *checkinRepository) GetCheckinsByUserIDAndDate(ctx context.Context, userID string, date string) ([]models.Checkin, error) {
	var checkins []models.Checkin

	// 查询数据库中符合条件的打卡记录
	query := `
		SELECT id, 
		       user_id, 
		       start_date, 
		       end_date, 
		       checkin_count, 
		       target_checkin_count, 
		       title, 
		       icon, 
		       motivation_message
		FROM checkins
		WHERE user_id = ? 
		  AND start_date <= ? 
		  AND end_date >= ?
	`

	// 执行查询，使用日期范围过滤
	rows, err := r.db.Raw(query, userID, date, date).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 逐行扫描结果
	for rows.Next() {
		var checkin models.Checkin
		var checkinCountJSON []byte // 用来存储 checkin_count 字段的字节数组

		if err := rows.Scan(&checkin.ID, &checkin.UserID, &checkin.StartDate, &checkin.EndDate, &checkinCountJSON, &checkin.TargetCheckinCount, &checkin.Title, &checkin.Icon, &checkin.MotivationMessage); err != nil {
			return nil, err
		}

		// 将 checkinCountJSON 字节数组解析为 map[string]int
		checkinCount := make(map[string]int)
		if err := json.Unmarshal(checkinCountJSON, &checkinCount); err != nil {
			return nil, err
		}

		// 将解析后的 checkinCount 存回 []byte
		checkinCountBytes, err := json.Marshal(checkinCount)
		if err != nil {
			return nil, err
		}

		// 将解析后的 checkinCount 赋值给 checkin
		checkin.CheckinCount = checkinCountBytes

		// 将 checkin 添加到结果列表
		checkins = append(checkins, checkin)
	}

	// 检查是否有错误发生
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkins, nil
}

// GetCheckinByID 获取指定ID的打卡记录
func (r *checkinRepository) GetCheckinByID(ctx context.Context, checkinID int) (*models.Checkin, error) {
	var checkin models.Checkin

	// 查询数据库中符合条件的打卡记录
	query := `
		SELECT id, 
		       user_id, 
		       start_date, 
		       end_date, 
		       checkin_count, 
		       target_checkin_count, 
		       title, 
		       icon, 
		       motivation_message
		FROM checkins
		WHERE id = ?
	`
	// 执行查询
	rows, err := r.db.Raw(query, checkinID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 如果查询不到记录，则返回 nil
	if !rows.Next() {
		return nil, fmt.Errorf("checkin with id %d not found", checkinID)
	}

	var checkinCountJSON []byte // 用来存储 checkin_count 字段的字节数组

	// 扫描结果到 checkin 结构体
	if err := rows.Scan(&checkin.ID, &checkin.UserID, &checkin.StartDate, &checkin.EndDate, &checkinCountJSON, &checkin.TargetCheckinCount, &checkin.Title, &checkin.Icon, &checkin.MotivationMessage); err != nil {
		return nil, err
	}

	// 将 checkinCountJSON 字节数组解析为 map[string]int
	checkinCount := make(map[string]int)
	if err := json.Unmarshal(checkinCountJSON, &checkinCount); err != nil {
		return nil, err
	}

	// 将解析后的 checkinCount 存回 []byte
	// 需要把 map[string]int 再转换成 []byte
	checkinCountBytes, err := json.Marshal(checkinCount)
	if err != nil {
		return nil, err
	}

	// 将解析后的 checkinCount 赋值给 checkin
	checkin.CheckinCount = checkinCountBytes

	// 返回查询到的单一打卡记录
	return &checkin, nil
}

func (r *checkinRepository) IncrementCheckinCount(ctx context.Context, checkinID int, date string) (*models.CheckinWithDecodedCount, error) {
	var checkin models.Checkin

	// 查询数据库中指定 checkin_id 的打卡记录
	query := `
        SELECT id, 
               user_id, 
               start_date, 
               end_date, 
               checkin_count, 
               target_checkin_count, 
               title, 
               icon, 
               motivation_message
        FROM checkins
        WHERE id = ?
    `
	// 使用 Raw 查询并将结果扫描到 checkin 结构体
	err := r.db.Raw(query, checkinID).Scan(&checkin).Error
	if err != nil {
		return nil, fmt.Errorf("checkin with id %d not found: %w", checkinID, err)
	}

	// checkin.CheckinCount 是 []byte 类型，因为数据库中存储的是 JSON
	// 解析 checkin_count 字段为 map[string]int
	checkinCount := make(map[string]int)
	if len(checkin.CheckinCount) > 0 {
		// 这里将 []byte 类型的 checkin_count 字段解码为 map[string]int
		err := json.Unmarshal(checkin.CheckinCount, &checkinCount)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal checkin_count: %w", err)
		}
	}

	// 判断当天的打卡次数是否已经达到目标次数
	if count, exists := checkinCount[date]; exists && count >= checkin.TargetCheckinCount {
		// 如果已完成打卡，则返回错误
		return nil, fmt.Errorf("checkin already completed for the date %s", date)
	}

	// 增加指定日期的打卡次数
	if count, exists := checkinCount[date]; exists {
		checkinCount[date] = count + 1
	} else {
		checkinCount[date] = 1
	}

	// 将更新后的 checkin_count 转回 JSON 格式
	updatedCheckinCountJSON, err := json.Marshal(checkinCount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated checkin_count: %w", err)
	}

	// 更新数据库中的 checkin_count 字段
	updateQuery := `
        UPDATE checkins
        SET checkin_count = ?
        WHERE id = ?
    `
	err = r.db.Exec(updateQuery, updatedCheckinCountJSON, checkinID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update checkin record with id %d: %w", checkinID, err)
	}

	// 将更新后的 checkin_count 返回到结构体
	checkin.CheckinCount = updatedCheckinCountJSON

	// 返回包含解码后的 Checkin 和解码后的 checkin_count 的结构体
	return &models.CheckinWithDecodedCount{
		Checkin:        &checkin,
		DecodedCheckin: checkinCount,
	}, nil
}

// UpdateCheckin 更新打卡记录到数据库
func (r *checkinRepository) UpdateCheckin(ctx context.Context, checkin *models.Checkin) error {
	// 更新数据库中的打卡记录
	query := `
		UPDATE checkins SET
			title = ?, 
			start_date = ?, 
			end_date = ?, 
			checkin_count = ?, 
			target_checkin_count = ?, 
			icon = ?, 
			motivation_message = ?
		WHERE id = ?
	`

	if err := r.db.Exec(query,
		checkin.Title,
		checkin.StartDate,
		checkin.EndDate,
		checkin.CheckinCount, // 这里传入 JSON 格式的 checkinCount
		checkin.TargetCheckinCount,
		checkin.Icon,
		checkin.MotivationMessage,
		checkin.ID,
	).Error; err != nil {
		return fmt.Errorf("failed to update checkin: %v", err)
	}

	return nil
}

// UpdateCount 更新打卡记录中的 checkin_count 字段
func (r *checkinRepository) UpdateCount(ctx context.Context, checkinID int, updatedCheckinCount []byte) error {
	// 更新数据库中的 checkin_count 字段
	query := `UPDATE checkins SET checkin_count = ? WHERE id = ?`

	if err := r.db.Exec(query, updatedCheckinCount, checkinID).Error; err != nil {
		return err
	}

	return nil
}

// DeleteCheckin 删除打卡记录
func (r *checkinRepository) DeleteCheckin(ctx context.Context, checkinID int) error {
	// 执行删除操作
	if err := r.db.Where("id = ?", checkinID).Delete(&models.Checkin{}).Error; err != nil {
		return err
	}
	return nil
}
