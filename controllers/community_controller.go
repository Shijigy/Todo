package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"ToDo/utils"
	"encoding/json"
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	var post models.CommunityPost

	// 解析 multipart/form-data 数据
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to parse form data"})
		return
	}

	// 从表单中获取 user_id 和 content 字段
	userID := r.FormValue("user_id")
	content := r.FormValue("content")

	// 确保 user_id 和 content 是有效的
	if userID == "" || content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "user_id and content are required"})
		return
	}

	// 设置 post 的 user_id 和 content
	post.UserID = userID
	post.Content = content

	// 检查是否上传了图片文件
	if file, _, err := r.FormFile("file"); err == nil && file != nil {
		// 上传图片到七牛云
		imageURL, err := utils.UploadImageToQiNiu(r)
		if err != nil {
			// 返回上传图片的错误
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.Response{Error: "Failed to upload image"})
			return
		}

		// 将图片链接赋值给 post.ImageURL
		post.ImageURL = imageURL
	}

	post.CommentCount = 0

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

// GetPosts 获取指定用户社区动态列表，支持分页、筛选、排序
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
	posts, err := service.GetCommunityPostsByUserIDService(r.Context(), service.Repo, page, limit, "", userID, sort)
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

// GetAllPosts 获取所有社区动态
func GetAllPosts(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	// 获取查询参数
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort") // 排序方式，默认为时间降序

	// 设置默认分页值
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	// 调用服务层获取所有社区帖子
	posts, err := service.GetAllCommunityPostsService(r.Context(), service.Repo, page, limit, sort)
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

// CheckLikeStatus 检查用户是否点赞动态
func CheckLikeStatus(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
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

	// 调用服务层的 CheckLikeStatusService 处理点赞状态
	status, err := service.CheckLikeStatusService(r.Context(), requestBody.UserID, requestBody.PostID)
	if err != nil {
		// 如果出错，返回错误信息
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回点赞状态
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Like status fetched successfully", Data: status})
}

// CreateComment 处理发布评论的请求
func CreateComment(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	var comment models.Comment

	// 解析 JSON 数据
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Failed to parse comment data"})
		return
	}

	// 调用服务层处理评论创建，并获取用户信息
	createdComment, username, avatarURL, err := service.CreateCommentService(r.Context(), comment)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回成功信息和详细数据
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Message: "Comment created successfully",
		Data: map[string]interface{}{
			"post_id":    createdComment.PostID,
			"user_id":    createdComment.UserID,
			"username":   username,
			"avatar_url": avatarURL,
			"content":    createdComment.Content,
			"created_at": createdComment.CreatedAt,
		},
	})
}

// DeleteComment 处理删除评论的请求
func DeleteComment(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	// 从 URL 查询参数中获取评论 ID
	commentID := r.URL.Query().Get("comment_id")
	if commentID == "" {
		// 如果没有传递 comment_id 参数，返回错误
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "comment_id is required"})
		return
	}

	// 调用服务层删除评论
	err := service.DeleteCommentService(r.Context(), commentID)
	if err != nil {
		// 如果服务层删除评论失败，返回 500 错误
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 如果删除成功，返回成功消息
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Comment deleted successfully"})
}

func GetComments(w http.ResponseWriter, r *http.Request, service services.CommunityService) {
	postID := r.URL.Query().Get("post_id")

	// 获取评论列表（包含用户信息）
	comments, err := service.GetCommentsService(r.Context(), postID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回评论列表
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Comments fetched successfully", Data: comments})
}
