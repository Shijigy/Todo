package repositories

import (
	"ToDo/models"
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// CommunityRepository 社区仓库接口
type CommunityRepository interface {
	CreatePost(ctx context.Context, post *models.CommunityPost) (*models.CommunityPost, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]models.CommunityPost, error)
	GetAllPosts(ctx context.Context) ([]models.CommunityPost, error)
	GetFilteredPosts(ctx context.Context, tags, userID, sort string, offset, limit int) ([]models.CommunityPost, error)
	IncrementLikesCount(ctx context.Context, postID string) error
	DecrementLikesCount(ctx context.Context, postID string) error
	DeletePost(ctx context.Context, id string) error
	GetLikesCount(ctx context.Context, postID string) (int, error)
}

// 社区仓库实现
type communityRepository struct {
	db *gorm.DB
}

// NewCommunityRepository 创建社区仓库实例
func NewCommunityRepository(db *gorm.DB) CommunityRepository {
	return &communityRepository{db: db}
}

// LikeRepository 点赞仓库接口
type LikeRepository interface {
	AddLike(ctx context.Context, userID, postID string) error
	CheckLikeExist(ctx context.Context, userID, postID string) (bool, error)
	RemoveLike(ctx context.Context, userID, postID string) error
}

// Like 仓库实现
type likeRepository struct {
	db *gorm.DB
}

// NewLikeRepository 创建点赞仓库实例
func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// CreatePost 创建社区动态
func (r *communityRepository) CreatePost(ctx context.Context, post *models.CommunityPost) (*models.CommunityPost, error) {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	if err := r.db.Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// GetPostsByUserID 根据用户 ID 获取社区动态
func (r *communityRepository) GetPostsByUserID(ctx context.Context, userID string) ([]models.CommunityPost, error) {
	var posts []models.CommunityPost
	if err := r.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetAllPosts 获取所有社区动态
func (r *communityRepository) GetAllPosts(ctx context.Context) ([]models.CommunityPost, error) {
	var posts []models.CommunityPost
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetFilteredPosts 获取符合条件的社区动态，支持分页、筛选、排序
func (r *communityRepository) GetFilteredPosts(ctx context.Context, tags, userID, sort string, offset, limit int) ([]models.CommunityPost, error) {
	var posts []models.CommunityPost
	query := r.db.Model(&models.CommunityPost{})

	// 标签筛选
	if tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	// 用户ID筛选
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 排序，默认按时间降序
	if sort == "" {
		sort = "desc"
	}
	if sort == "asc" {
		query = query.Order("created_at asc")
	} else {
		query = query.Order("created_at desc")
	}

	// 分页
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// DeletePost 删除社区动态
func (r *communityRepository) DeletePost(ctx context.Context, id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.CommunityPost{}).Error; err != nil {
		return err
	}
	return nil
}

// IncrementLikesCount 增加社区动态的点赞数
func (r *communityRepository) IncrementLikesCount(ctx context.Context, postID string) error {
	var post models.CommunityPost
	if err := r.db.Model(&post).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// AddLike 添加点赞记录
func (r *likeRepository) AddLike(ctx context.Context, userID, postID string) error {
	like := models.Like{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}
	if err := r.db.Create(&like).Error; err != nil {
		return err
	}
	return nil
}

// CheckLikeExist 检查是否已经点赞
func (r *likeRepository) CheckLikeExist(ctx context.Context, userID, postID string) (bool, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return err == nil, nil
}

// DecrementLikesCount 减少社区动态的点赞数
func (r *communityRepository) DecrementLikesCount(ctx context.Context, postID string) error {
	var post models.CommunityPost
	err := r.db.Model(&post).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

// RemoveLike 删除用户对帖子点赞记录
func (r *likeRepository) RemoveLike(ctx context.Context, userID, postID string) error {
	// 使用GORM删除对应的点赞记录
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetLikesCount 获取特定帖子点赞数
func (r *communityRepository) GetLikesCount(ctx context.Context, postID string) (int, error) {
	var count int
	err := r.db.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
