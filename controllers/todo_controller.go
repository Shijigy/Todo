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
	if todo.UpdatedAt.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "UpdatedAt is required"})
		return
	}

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

func GetTodos(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	// 获取查询参数
	userID := r.URL.Query().Get("user_id")
	updatedAtStr := r.URL.Query().Get("updated_at") // 接收更新日期

	// 解析 updated_at 字符串为 time.Time，并只获取到日期部分
	var updatedAt time.Time
	var err error
	if updatedAtStr != "" {
		updatedAt, err = time.Parse("2006-01-02", updatedAtStr) // 仅解析到日期，忽略时间部分
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.Response{Error: "无效的 updated_at 格式，预期为 YYYY-MM-DD"})
			return
		}
		// 设置时间为00:00:00，确保只比对日期部分
		updatedAt = updatedAt.Add(time.Hour * 24 * 0) // 将时间设置为午夜（00:00:00），忽略时间
	}

	// 传递 context 到服务层
	ctx := r.Context()
	todos, err := todoService.GetTodosService(ctx, userID, updatedAt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "任务获取成功", Data: todos})
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
		Title string `json:"title"`
	}

	// 使用 ShouldBind 来解析请求体
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid input"})
		return
	}

	// 传递 context 到服务层
	ctx := c.Request.Context()
	updatedTodo, err := todoService.MarkTodoAsCompletedService(ctx, request.Title)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Error: err.Error()})
		return
	}

	// 直接使用 c.JSON 设置响应头、状态码以及响应数据
	c.JSON(http.StatusOK, models.Response{Message: "Todo marked as completed", Data: updatedTodo})
}
