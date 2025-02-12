package models

import "time"

// CommunityPost 社区动态模型
type CommunityPost struct {
	ID         string    `json:"id" bson:"_id"`
	UserID     string    `json:"user_id" bson:"user_id"`
	Content    string    `json:"content" bson:"content"`
	ImageURL   string    `json:"image_url" bson:"image_url"` // 动态附带图片
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
	Tags       string    `json:"tags" bson:"tags"`               // 动态的标签
	LikesCount int       `json:"likes_count" bson:"likes_count"` // 点赞数
}

// NewCommunityPost 创建社区动态
func NewCommunityPost(userID, content, imageURL string) *CommunityPost {
	return &CommunityPost{
		UserID:    userID,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Like 点赞模型
type Like struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	PostID    string    `json:"post_id" bson:"post_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
