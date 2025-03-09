package models

import "time"

// CommunityPost 社区动态模型
type CommunityPost struct {
	ID           int       `json:"id" bson:"_id"` // 使用 int 类型，匹配数据库中的自增 ID
	UserID       string    `json:"user_id" bson:"user_id"`
	Content      string    `json:"content" bson:"content"`
	ImageURL     string    `json:"image_url" bson:"image_url"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	Tags         string    `json:"tags" bson:"tags"`                   // 标签
	LikesCount   int       `json:"likes_count" bson:"likes_count"`     // 点赞数
	CommentCount int       `json:"comment_count" bson:"comment_count"` // 评论数量字段
}

// Like 点赞模型
type Like struct {
	ID        int       `json:"id" bson:"_id"` // 点赞记录 ID
	UserID    string    `json:"user_id" bson:"user_id"`
	PostID    string    `json:"post_id" bson:"post_id"` // 关联的社区动态 ID
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// Comment 评论模型
type Comment struct {
	ID        int       `json:"id"`         // 评论的唯一标识符
	PostID    int       `json:"post_id"`    // 关联的动态 ID
	UserID    string    `json:"user_id"`    // 评论的用户 ID
	Content   string    `json:"content"`    // 评论内容
	CreatedAt time.Time `json:"created_at"` // 评论时间
}

// CommentWithUser 包含评论和用户信息的结构体
type CommentWithUser struct {
	Comment Comment `json:"comment"`
	User    *User   `json:"user"`
}
