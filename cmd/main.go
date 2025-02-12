package main

import (
	"ToDo/config"
	"ToDo/controllers"
	"ToDo/middlewares"
	"ToDo/repositories"
	"ToDo/services"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var db *gorm.DB
var userRepo repositories.UserRepository
var userService services.UserService
var checkinRepo repositories.CheckinRepository
var checkinService services.CheckinService
var communityRepo repositories.CommunityRepository
var communityService services.CommunityService
var todoRepo repositories.TodoRepository
var todoService services.TodoService

// 启动入口
func main() {
	// 加载配置
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// 输出加载的配置
	fmt.Println("Server address:", config.ServerAddress)
	fmt.Println("Database host:", config.Database.Host)

	// 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DbName,
	)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// 初始化各个仓库和服务
	userRepo = repositories.NewUserRepository(db)
	userService = services.NewUserService(userRepo)
	checkinRepo = repositories.NewCheckinRepository(db)
	checkinService = services.NewCheckinService(checkinRepo)

	// 初始化社区仓库和服务
	communityRepo = repositories.NewCommunityRepository(db)
	communityService = services.NewCommunityService(communityRepo)

	// 初始化待办任务仓库和服务
	todoRepo = repositories.NewTodoRepository(db)
	todoService = services.NewTodoService(todoRepo)

	// 初始化 Gin 路由器
	r := gin.Default()

	// 使用日志中间件
	r.Use(middlewares.LoggingMiddleware)

	// CORS 配置
	r.Use(cors.Default())
	// 路由设置
	r.POST("/auth/login", func(c *gin.Context) {
		controllers.Login(c.Writer, c.Request, userService)
	})
	r.POST("/auth/register", func(c *gin.Context) {
		emailService := services.NewEmailService()
		controllers.Register(c.Writer, c.Request, userService, emailService)
	})

	// Todo 路由
	r.POST("/todo", func(c *gin.Context) {
		controllers.CreateTodo(c.Writer, c.Request, todoService)
	})
	r.GET("/todo/:id", func(c *gin.Context) {
		controllers.GetTodo(c.Writer, c.Request, todoService)
	})
	r.PUT("/todo/:id", func(c *gin.Context) {
		controllers.UpdateTodo(c.Writer, c.Request, todoService)
	})
	r.DELETE("/todo/:id", func(c *gin.Context) {
		controllers.DeleteTodo(c.Writer, c.Request, todoService)
	})

	// Checkin 路由
	r.POST("/checkin", func(c *gin.Context) {
		controllers.Checkin(c.Writer, c.Request, checkinService)
	})

	// 社区动态路由
	r.POST("/community/post", func(c *gin.Context) {
		controllers.CreatePost(c.Writer, c.Request, communityService)
	})
	r.GET("/community/list", func(c *gin.Context) {
		controllers.GetPosts(c.Writer, c.Request, communityService)
	})
	r.POST("/community/:post_id/like", func(c *gin.Context) {
		// 提取 path 参数 (post_id)
		postID := c.Param("post_id")

		// 解析请求体中的 UserID
		var requestBody struct {
			UserID string `json:"user_id"` // 修改字段为小写 user_id
		}

		// 解析请求体
		if err := json.NewDecoder(c.Request.Body).Decode(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		// 校验 UserID 是否有效
		if requestBody.UserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is required"})
			return
		}

		// 调用 LikePost 控制器并传入 postID 和 UserID
		controllers.LikePost(c.Writer, c.Request, communityService, postID, requestBody.UserID)
	})

	// 启动 HTTP 服务
	fmt.Println("Server running on", config.ServerAddress)

	// 启动 HTTP 服务并监听端口
	if err := r.Run(config.ServerAddress); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
