package models

import "time"

// Checkin 打卡任务模型
type Checkin struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	CheckinAt time.Time `json:"checkin_at" bson:"checkin_at"`
	Status    string    `json:"status" bson:"status"` // 状态: successful, failed
}
