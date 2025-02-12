package models

import "time"

// CommunityPost 社区动态模型
type CommunityPost struct {
	ID         int       `json:"id" bson:"_id"` // 使用 int 类型，匹配数据库中的自增 ID
	UserID     string    `json:"user_id" bson:"user_id"`
	Content    string    `json:"content" bson:"content"`
	ImageURL   string    `json:"image_url" bson:"image_url"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
	Tags       string    `json:"tags" bson:"tags"`               // 标签
	LikesCount int       `json:"likes_count" bson:"likes_count"` // 点赞数
}

// NewCommunityPost 创建社区动态
func NewCommunityPost(userID, content, imageURL, tags string) *CommunityPost {
	return &CommunityPost{
		UserID:    userID,
		Content:   content,
		ImageURL:  imageURL,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Like 点赞模型
type Like struct {
	ID        int       `json:"id" bson:"_id"` // 点赞记录 ID
	UserID    string    `json:"user_id" bson:"user_id"`
	PostID    string    `json:"post_id" bson:"post_id"` // 关联的社区动态 ID
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
