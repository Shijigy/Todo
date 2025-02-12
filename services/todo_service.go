package services

import (
	"ToDo/models"
	"ToDo/repositories"
	"context"
	"errors"
	"sync"
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

// GetTodoService 根据标题获取待办任务
func (s *TodoService) GetTodoService(ctx context.Context, title string) (*models.Todo, error) {
	if title == "" {
		return nil, errors.New("todo title cannot be empty")
	}

	todo, err := s.Repo.GetTodoByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// UpdateTodoStatusService 更新任务状态
func (s *TodoService) UpdateTodoStatusService(ctx context.Context, title string, todo models.Todo) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	// 获取现有任务
	existingTodo, err := s.Repo.GetTodoByTitle(ctx, title)
	if err != nil {
		return errors.New("todo not found")
	}

	// 更新任务状态
	existingTodo.Status = todo.Status
	return s.Repo.UpdateTodoStatus(ctx, existingTodo)
}

// DeleteTodoService 删除任务
func (s *TodoService) DeleteTodoService(ctx context.Context, title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	// 删除任务
	return s.Repo.DeleteTodoByTitle(ctx, title)
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
