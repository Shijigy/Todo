package models

type Checkin struct {
	ID        int    `json:"id" bson:"_id"`
	UserID    string `json:"user_id" bson:"user_id"`
	StartDate string `json:"start_date" bson:"start_date"`
	EndDate   string `json:"end_date" bson:"end_date"`
	// CheckinCount       map[string]int json:"checkin_count" bson:"checkin_count" // 每天的打卡次数
	CheckinCount       []byte `json:"checkin_count"`
	TargetCheckinCount int    `json:"target_checkin_count" bson:"target_checkin_count"`
	Title              string `json:"title" bson:"title"`
	Icon               int    `json:"icon" bson:"icon"`
	MotivationMessage  string `json:"motivation_message" bson:"motivation_message"`
}

// CheckinWithDecodedCount 用于包含原始的 Checkin 和解码后的 CheckinCount
type CheckinWithDecodedCount struct {
	Checkin        *Checkin       `json:"checkin"`
	DecodedCheckin map[string]int `json:"decoded_checkin_count"`
}

// CheckinCompletionRequest 请求参数结构体
type CheckinCompletionRequest struct {
	CheckinID int    `json:"checkin_id"`
	Date      string `json:"date"`
}
