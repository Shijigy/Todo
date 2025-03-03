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
	GetTodoByID(ctx context.Context, id string) (*models.Todo, error)
	UpdateTodoStatus(ctx context.Context, todo *models.Todo) error
	DeleteTodoByID(ctx context.Context, id string) error
	GetAllTodos(ctx context.Context) ([]models.Todo, error)
	GetTodoByTitle(ctx context.Context, title string) (*models.Todo, error)
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
	if todo.UpdatedAt.IsZero() {
		todo.UpdatedAt = time.Now() // 默认使用当前时间
	}

	// 保存 todo 到数据库
	if err := r.db.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

// GetAllTodos 获取所有待办任务
func (r *todoRepository) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	var todos []models.Todo
	// 获取所有任务
	if err := r.db.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

// GetTodoByID 根据任务ID获取待办任务
func (r *todoRepository) GetTodoByID(ctx context.Context, id string) (*models.Todo, error) {
	var todo models.Todo
	// 根据ID查找任务
	if err := r.db.Where("id = ?", id).First(&todo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil // 返回 nil 表示任务没有找到
		}
		return nil, err
	}
	return &todo, nil
}

// GetTodoByTitle 根据任务标题获取待办任务
func (r *todoRepository) GetTodoByTitle(ctx context.Context, title string) (*models.Todo, error) {
	var todo models.Todo
	// 根据标题查找任务
	if err := r.db.Where("title = ?", title).First(&todo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil // 返回 nil 表示任务没有找到
		}
		return nil, err
	}
	return &todo, nil
}

// UpdateTodoStatus 更新待办任务状态
func (r *todoRepository) UpdateTodoStatus(ctx context.Context, todo *models.Todo) error {
	todo.UpdatedAt = time.Now()
	// 更新任务
	if err := r.db.Save(todo).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTodoByID 根据任务ID删除待办任务
func (r *todoRepository) DeleteTodoByID(ctx context.Context, id string) error {
	// 根据ID删除任务
	if err := r.db.Where("id = ?", id).Delete(&models.Todo{}).Error; err != nil {
		return err
	}
	return nil
}
