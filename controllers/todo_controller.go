package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// CreateTodo 创建待办任务
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

	// 离线模式
	if isOffline {
		todo.CreatedAt = time.Now()
		todo.UpdatedAt = time.Now()

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

// GetTodo 获取单个待办任务
func GetTodo(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	params := mux.Vars(r)
	todoID := params["id"]

	// 传递 context 到服务层
	ctx := r.Context()
	todo, err := todoService.GetTodoService(ctx, todoID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Todo fetched successfully", Data: todo})
}

// UpdateTodo 更新待办任务
func UpdateTodo(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	params := mux.Vars(r)
	todoID := params["id"]

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 传递 context 到服务层
	ctx := r.Context()
	err := todoService.UpdateTodoStatusService(ctx, todoID, todo)
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
func DeleteTodo(w http.ResponseWriter, r *http.Request, todoService services.TodoService) {
	params := mux.Vars(r)
	todoID := params["id"]

	// 传递 context 到服务层
	ctx := r.Context()
	err := todoService.DeleteTodoService(ctx, todoID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Todo deleted successfully"})
}
