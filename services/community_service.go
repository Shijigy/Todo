package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"errors"
	"strconv"
)

type CommunityService struct {
	Repo     repositories.CommunityRepository
	LikeRepo repositories.LikeRepository
}

// NewCommunityService 是一个构造函数，返回一个具体的 CommunityService 实例
func NewCommunityService(repo repositories.CommunityRepository) CommunityService {
	return CommunityService{Repo: repo}
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

	return posts, nil
}

// GetCommunityPostsByUserIDService 获取指定用户的社区动态
func (s CommunityService) GetCommunityPostsByUserIDService(ctx context.Context, userID string, repo repositories.CommunityRepository) ([]models.CommunityPost, error) {
	posts, err := repo.GetPostsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("no posts found for the user")
	}
	return posts, nil
}

// LikePostService 点赞社区动态
func (s CommunityService) LikePostService(ctx context.Context, userID, postID string) error {
	// 验证用户是否已经点赞
	likeExists, err := s.LikeRepo.CheckLikeExist(ctx, userID, postID)
	if err != nil {
		return err
	}

	if likeExists {
		return errors.New("you have already liked this post")
	}

	// 添加点赞记录
	err = s.LikeRepo.AddLike(ctx, userID, postID)
	if err != nil {
		return err
	}

	// 更新动态的点赞数
	err = s.Repo.IncrementLikesCount(ctx, postID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCommunityPostService 删除社区动态
func (s CommunityService) DeleteCommunityPostService(ctx context.Context, id string, repo repositories.CommunityRepository) error {
	err := repo.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
