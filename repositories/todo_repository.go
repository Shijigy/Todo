package repositories

import (
	"ToDo/models"
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// TodoRepository 待办任务仓库接口
type TodoRepository interface {
	CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)
	GetTodoByTitle(ctx context.Context, title string) (*models.Todo, error)
	UpdateTodoStatus(ctx context.Context, todo *models.Todo) error
	DeleteTodoByTitle(ctx context.Context, title string) error
}

// 待办任务仓库实现
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 创建待办任务仓库实例
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

// CreateTodo 创建待办任务
func (r *todoRepository) CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error) {
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	if err := r.db.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

// GetTodoByTitle 根据任务标题获取待办任务
func (r *todoRepository) GetTodoByTitle(ctx context.Context, title string) (*models.Todo, error) {
	var todo models.Todo
	if err := r.db.Where("title = ?", title).First(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodoStatus 更新待办任务状态
func (r *todoRepository) UpdateTodoStatus(ctx context.Context, todo *models.Todo) error {
	todo.UpdatedAt = time.Now()
	if err := r.db.Save(todo).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTodoByTitle 根据标题删除待办任务
func (r *todoRepository) DeleteTodoByTitle(ctx context.Context, title string) error {
	if err := r.db.Where("title = ?", title).Delete(&models.Todo{}).Error; err != nil {
		return err
	}
	return nil
}
