package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateTodo(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	// 检查是否为离线模式
	isOffline := r.Header.Get("Is-Offline") == "true"

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 解析 updated_at 字符串为 time.Time
	if todo.UpdatedAt == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "UpdatedAt is required"})
		return
	}

	// 如果 UpdatedAt 字符串为空，可以设置为当前时间
	parsedUpdatedAt, err := time.Parse("2006-01-02", todo.UpdatedAt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid updated_at format, expected YYYY-MM-DD"})
		return
	}

	// 将 parsedUpdatedAt 转换回字符串类型
	todo.UpdatedAt = parsedUpdatedAt.Format("2006-01-02")

	// 离线模式
	if isOffline {
		todo.CreatedAt = time.Now()

		// 在离线模式下，传递 context.Background() 而不是 r.Context()
		err := todoService.SaveOffline(context.Background(), todo)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.Response{Error: "Failed to save todo offline"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.Response{Message: "Todo created successfully in offline mode", Data: todo})
		return
	}

	// 传递 context 到服务层
	ctx := r.Context()
	createdTodo, err := todoService.CreateTodoService(ctx, todo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{Message: "Todo created successfully", Data: createdTodo})
}

// GetTodos 获取指定用户和日期范围内的任务
func GetTodos(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	// 从请求体中解析 JSON 数据
	var requestData struct {
		UserID    string `json:"user_id"`
		UpdatedAt string `json:"updated_at"` // 接收更新日期
	}

	// 解析请求体
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid request body"})
		return
	}

	// 解析 updated_at 字符串为 time.Time，并只获取到日期部分
	var updatedAt time.Time
	if requestData.UpdatedAt != "" {
		updatedAt, err = time.Parse("2006-01-02", requestData.UpdatedAt) // 解析到日期
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.Response{Error: "Invalid updated_at format, expected YYYY-MM-DD"})
			return
		}
		// 设置时间为00:00:00，确保只比对日期部分
		updatedAt = updatedAt.Add(time.Hour * 24 * 0) // 将时间设置为午夜（00:00:00），忽略时间
	}

	// 将 updatedAt 转换为字符串
	updatedAtFormatted := updatedAt.Format("2006-01-02")

	// 传递 context 到服务层
	ctx := r.Context()
	todos, err := todoService.GetTodosService(ctx, requestData.UserID, updatedAtFormatted)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回任务列表
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

// UpdateTodo 更新待办任务
func UpdateTodo(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 传递 context 到服务层
	ctx := r.Context()
	err := todoService.UpdateTodoStatusService(ctx, todo.ID, todo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Todo updated successfully"})
}

// DeleteTodo 删除待办任务
func DeleteTodo(c *gin.Context, todoService services.TodoService) {
	todoID := c.Param("id")

	// 传递 context 到服务层
	ctx := c.Request.Context()
	err := todoService.DeleteTodoService(ctx, todoID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "Todo deleted successfully"})
}

// MarkTodoAsCompleted 标记任务为已完成
func MarkTodoAsCompleted(c *gin.Context, todoService services.TodoService) {
	var request struct {
		ID string `json:"id"` // 根据 ID 查找任务
	}

	// 使用 ShouldBind 来解析请求体
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid input"})
		return
	}

	// 传递 context 到服务层
	ctx := c.Request.Context()
	updatedTodo, err := todoService.MarkTodoAsCompletedService(ctx, request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Error: err.Error()})
		return
	}

	// 直接使用 c.JSON 设置响应头、状态码以及响应数据
	c.JSON(http.StatusOK, models.Response{Message: "Todo marked as completed", Data: updatedTodo})
}
