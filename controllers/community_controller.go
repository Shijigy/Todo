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
	tags := r.URL.Query().Get("tags")      // 标签筛选
	userID := r.URL.Query().Get("user_id") // 用户 ID 筛选
	sort := r.URL.Query().Get("sort")      // 排序方式，默认为时间降序

	// 设置默认分页值
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	// 调用服务层获取符合条件的社区帖子
	posts, err := service.GetCommunityPostsService(r.Context(), service.Repo, page, limit, tags, userID, sort)
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

// LikePost 点赞社区动态
func LikePost(w http.ResponseWriter, r *http.Request, service services.CommunityService, postID string, userID string) {
	// 校验 UserID 和 PostID 是否有效
	if userID == "" || postID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "UserID and PostID are required"})
		return
	}

	// 调用服务层的 LikePostService 处理点赞逻辑
	ctx := r.Context() // 获取上下文
	err := service.LikePostService(ctx, userID, postID)
	if err != nil {
		// 如果出错，返回错误信息
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

}
