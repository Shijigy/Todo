package models

import "time"

// Checkin 打卡任务模型
type Checkin struct {
	ID                 string    `json:"id" bson:"_id"`
	UserID             string    `json:"user_id" bson:"user_id"`
	CheckinAt          time.Time `json:"checkin_at" bson:"checkin_at"`
	Status             string    `json:"status" bson:"status"`                             // 状态: successful, failed
	CheckinCount       int       `json:"checkin_count" bson:"checkin_count"`               // 当前已完成的打卡次数
	TargetCheckinCount int       `json:"target_checkin_count" bson:"target_checkin_count"` // 目标打卡次数
	Title              string    `json:"title" bson:"title"`                               // 打卡标题
}
