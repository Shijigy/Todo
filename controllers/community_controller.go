package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"encoding/json"
	"net/http"
)

// CreatePost 创建社区动态
func CreatePost(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	var post models.CommunityPost

	// 解析请求体中的数据
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 调用服务层创建社区帖子
	createdPost, err := service.CreateCommunityPostService(r.Context(), post, service.Repo)
	if err != nil {
		// 返回服务层错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功创建社区帖子，返回 201 状态码和响应信息
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{Message: "Post created successfully", Data: createdPost})
}

// GetPosts 获取社区动态列表，支持分页、筛选、排序
func GetPosts(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	// 获取查询参数
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")      // 排序方式，默认为时间降序
	userID := r.URL.Query().Get("user_id") // 根据用户 ID 获取动态

	// 如果没有提供 user_id，则返回错误
	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "user_id is required"})
		return
	}

	// 设置默认分页值
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	// 调用服务层获取符合条件的社区帖子，仅根据 user_id 进行筛选
	posts, err := service.GetCommunityPostsService(r.Context(), service.Repo, page, limit, "", userID, sort)
	if err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功获取社区帖子，返回 200 状态码和响应数据
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Posts retrieved successfully", Data: posts})
}

// DeletePost 删除社区动态
func DeletePost(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	// 获取动态 ID
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Post ID is required"})
		return
	}

	// 调用服务层删除社区动态
	err := service.DeleteCommunityPostService(r.Context(), id, service.Repo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功删除社区动态，返回 200 状态码
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Post deleted successfully"})
}

// LikePost 点赞社区动态
func LikePost(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	var requestBody struct {
		UserID string `json:"user_id"`
		PostID string `json:"post_id"`
	}

	// 解析请求体中的数据
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 校验 UserID 和 PostID 是否有效
	if requestBody.UserID == "" || requestBody.PostID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "UserID and PostID are required"})
		return
	}

	// 调用服务层的 LikePostService 处理点赞逻辑
	err := service.LikePostService(r.Context(), requestBody.UserID, requestBody.PostID)
	if err != nil {
		// 如果出错，返回错误信息
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功处理点赞
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Post liked successfully"})
}

// CancelLikePost 取消点赞社区动态
func CancelLikePost(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	var requestBody struct {
		UserID string `json:"user_id"`
		PostID string `json:"post_id"`
	}

	// 解析请求体中的数据
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 校验 UserID 和 PostID 是否有效
	if requestBody.UserID == "" || requestBody.PostID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "UserID and PostID are required"})
		return
	}

	// 调用服务层的 CancelLikePostService 处理取消点赞逻辑
	err := service.CancelLikePostService(r.Context(), requestBody.UserID, requestBody.PostID)
	if err != nil {
		// 如果出错，返回错误信息
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功处理取消点赞
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Post unliked successfully"})
}

// GetLikesCount 获取特定帖子的点赞数
func GetLikesCount(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	postID := r.URL.Query().Get("post_id")

	// 校验 PostID 是否有效
	if postID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "PostID is required"})
		return
	}

	// 调用服务层的 GetLikesCountService 获取点赞数
	likesCount, err := service.GetLikesCountService(r.Context(), postID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 成功返回点赞数
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Likes count retrieved successfully", Data: map[string]int{"likes_count": likesCount}})
}
