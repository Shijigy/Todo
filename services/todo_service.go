package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"errors"
	"sync"
	"time"
)

type TodoService struct {
	Repo         repositories.TodoRepository
	OfflineTodos map[string]models.Todo // 离线任务存储
	mu           sync.Mutex             // 用于并发控制
}

// NewTodoService 创建待办任务
func NewTodoService(repo repositories.TodoRepository) TodoService {
	return TodoService{
		Repo:         repo,
		OfflineTodos: make(map[string]models.Todo), // 初始化离线任务存储
	}
}

// SaveOffline 保存离线任务
func (s *TodoService) SaveOffline(ctx context.Context, todo models.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 保存到内存中的离线任务
	s.OfflineTodos[todo.ID] = todo
	return nil
}

// CreateTodoService 创建待办任务
func (s *TodoService) CreateTodoService(ctx context.Context, todo models.Todo) (*models.Todo, error) {
	// 校验待办任务的输入数据
	if todo.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	// 调用仓库层创建待办任务
	createdTodo, err := s.Repo.CreateTodo(ctx, &todo)
	if err != nil {
		return nil, err
	}
	return createdTodo, nil
}

func (s *TodoService) GetTodosService(ctx context.Context, userID string, updatedAt time.Time) ([]models.Todo, error) {
	// 获取所有任务
	todos, err := s.Repo.GetAllTodos(ctx)
	if err != nil {
		return nil, err
	}

	// 根据 user_id 和 updated_at 过滤任务
	var filteredTodos []models.Todo
	for _, todo := range todos {
		if todo.UserID == userID && (updatedAt.IsZero() || todo.UpdatedAt.After(updatedAt)) {
			filteredTodos = append(filteredTodos, todo)
		}
	}

	return filteredTodos, nil
}

// GetTodoService 根据标题获取待办任务
func (s *TodoService) GetTodoService(ctx context.Context, tittle string) (*models.Todo, error) {
	if tittle == "" {
		return nil, errors.New("todo title cannot be empty")
	}

	todo, err := s.Repo.GetTodoByID(ctx, tittle)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// UpdateTodoStatusService 更新任务状态
func (s *TodoService) UpdateTodoStatusService(ctx context.Context, id string, todo models.Todo) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	// 获取现有任务
	existingTodo, err := s.Repo.GetTodoByID(ctx, id) // 根据ID获取任务
	if err != nil {
		return errors.New("todo not found")
	}
	// 检查 existingTodo 是否为 nil
	if existingTodo == nil {
		return errors.New("todo not found, existingTodo is nil")
	}

	// 如果传入了 updated_at 字段，则保留传入的值，而不是修改为当前时间
	if !todo.UpdatedAt.IsZero() {
		existingTodo.UpdatedAt = todo.UpdatedAt // 保持传入的 updated_at 时间
	}

	// 更新其他字段
	if todo.Status != "" {
		existingTodo.Status = todo.Status
	}
	if todo.Title != "" {
		existingTodo.Title = todo.Title
	}
	if todo.Description != "" {
		existingTodo.Description = todo.Description
	}

	// 保存更新后的任务
	return s.Repo.UpdateTodoStatus(ctx, existingTodo)
}

// DeleteTodoService 删除任务
func (s *TodoService) DeleteTodoService(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	// 删除任务
	return s.Repo.DeleteTodoByID(ctx, id) // 根据ID删除任务
}

// SyncOfflineTodos 同步离线任务到数据库
func (s *TodoService) SyncOfflineTodos() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, todo := range s.OfflineTodos {
		// 将离线任务同步到数据库
		_, err := s.Repo.CreateTodo(context.Background(), &todo)
		if err != nil {
			return err
		}
	}

	// 清空离线任务
	s.OfflineTodos = make(map[string]models.Todo)
	return nil
}

// MarkTodoAsCompletedService 标记任务为已完成
func (s *TodoService) MarkTodoAsCompletedService(ctx context.Context, id string) (*models.Todo, error) {
	// 根据任务的 ID 查找任务
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}

	// 获取现有任务
	existingTodo, err := s.Repo.GetTodoByID(ctx, id) // 根据 ID 获取任务
	if err != nil {
		return nil, errors.New("todo not found")
	}

	// 更新任务状态为 completed
	existingTodo.Status = "completed"
	existingTodo.UpdatedAt = time.Now() // 更新时间戳

	// 保存更新后的任务
	err = s.Repo.UpdateTodoStatus(ctx, existingTodo)
	if err != nil {
		return nil, err
	}

	return existingTodo, nil
}
