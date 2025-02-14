package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"errors"
	"fmt"
	"strconv"
)

type CommunityService struct {
	Repo     repositories.CommunityRepository
	LikeRepo repositories.LikeRepository
}

// NewCommunityService 是一个构造函数，返回一个具体的 CommunityService 实例
func NewCommunityService(repo repositories.CommunityRepository, likeRepo repositories.LikeRepository) CommunityService {
	return CommunityService{
		Repo:     repo,
		LikeRepo: likeRepo,
	}
}

// CreateCommunityPostService 创建社区动态
func (s CommunityService) CreateCommunityPostService(ctx context.Context, post models.CommunityPost, repo repositories.CommunityRepository) (*models.CommunityPost, error) {
	createdPost, err := repo.CreatePost(ctx, &post)
	if err != nil {
		return nil, err
	}
	return createdPost, nil
}

// GetCommunityPostsService 获取社区动态列表，支持分页、筛选、排序
func (s CommunityService) GetCommunityPostsService(ctx context.Context, repo repositories.CommunityRepository, page, limit, tags, userID, sort string) ([]models.CommunityPost, error) {
	// 处理分页
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt

	// 调用仓库方法获取符合条件的社区帖子
	posts, err := repo.GetFilteredPosts(ctx, tags, userID, sort, offset, limitInt)
	if err != nil {
		return nil, err
	}

	// 返回包含动态内容、用户信息、发布时间、图片、标签等信息的帖子
	var resultPosts []models.CommunityPost
	for _, post := range posts {
		// 你可以在这里加入任何额外的逻辑来填充用户信息等内容，假设用户信息已包含在 `CommunityPost` 模型内
		resultPosts = append(resultPosts, post)
	}

	return resultPosts, nil
}

// GetCommunityPostsByUserIDService 获取指定用户的社区动态
func (s CommunityService) GetCommunityPostsByUserIDService(ctx context.Context, userID string, repo repositories.CommunityRepository) ([]models.CommunityPost, error) {
	posts, err := repo.GetPostsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("no posts found for the user")
	}
	return posts, nil
}

// DeleteCommunityPostService 删除社区动态
func (s CommunityService) DeleteCommunityPostService(ctx context.Context, id string, repo repositories.CommunityRepository) error {
	// 调用仓库层删除动态
	err := repo.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// LikePostService 点赞社区动态服务
func (s *CommunityService) LikePostService(ctx context.Context, userID, postID string) error {
	// 检查是否已经点赞
	liked, err := s.LikeRepo.CheckLikeExist(ctx, userID, postID)
	if err != nil {
		return err
	}

	if liked {
		return fmt.Errorf("user has already liked this post")
	}

	// 添加点赞记录
	err = s.LikeRepo.AddLike(ctx, userID, postID)
	if err != nil {
		return err
	}

	// 增加帖子点赞数
	err = s.Repo.IncrementLikesCount(ctx, postID)
	if err != nil {
		return err
	}

	return nil
}

// CancelLikePostService 取消点赞服务
func (s *CommunityService) CancelLikePostService(ctx context.Context, userID, postID string) error {
	// 验证用户是否已点赞
	liked, err := s.LikeRepo.CheckLikeExist(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check like existence: %v", err)
	}
	if !liked {
		return errors.New("cannot unlike a post that was not liked")
	}

	// 删除用户的点赞记录
	err = s.LikeRepo.RemoveLike(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to remove like: %v", err)
	}

	// 更新帖子点赞数
	err = s.Repo.DecrementLikesCount(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to decrement likes count: %v", err)
	}

	return nil
}

// GetLikesCountService 获取特定帖子的点赞数
func (s *CommunityService) GetLikesCountService(ctx context.Context, postID string) (int, error) {
	likesCount, err := s.Repo.GetLikesCount(ctx, postID)
	if err != nil {
		return 0, err
	}
	return likesCount, nil
}
