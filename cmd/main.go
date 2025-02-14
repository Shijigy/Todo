package main

import (
	"ToDo/config"
	"ToDo/controllers"
	"ToDo/middlewares"
	"ToDo/models"
	"ToDo/repositories"
	"ToDo/services"
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
var likeRepo repositories.LikeRepository // 确保声明 likeRepo
var emailService services.EmailService   // 声明 emailService

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
	emailService = services.NewEmailService(config.Email.SMTPServer, config.Email.FromEmail, config.Email.Password)
	userService = services.NewUserService(userRepo, emailService)
	checkinRepo = repositories.NewCheckinRepository(db)
	checkinService = services.NewCheckinService(checkinRepo)

	// 初始化服务
	communityRepo = repositories.NewCommunityRepository(db)
	likeRepo = repositories.NewLikeRepository(db)
	communityService = services.NewCommunityService(communityRepo, likeRepo)

	// 初始化待办任务仓库和服务
	todoRepo = repositories.NewTodoRepository(db)
	todoService = services.NewTodoService(todoRepo)

	// 初始化 Gin 路由器
	r := gin.Default()

	// 使用日志中间件
	r.Use(middlewares.LoggingMiddleware)

	// CORS 配置
	r.Use(cors.Default())

	r.POST("/auth/register", func(c *gin.Context) {
		// 解析请求体
		var req struct {
			User         models.User `json:"user"`
			CaptchaInput string      `json:"captchaInput"`
		}

		// 绑定请求体中的 JSON 数据到 req 结构体
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid input"})
			return
		}

		// 继续调用原来的 Register 函数，适配 c.Writer 和 c.Request
		controllers.Register(c.Writer, c.Request, userService, emailService)
	})

	r.POST("/auth/login", func(c *gin.Context) {
		controllers.Login(c.Writer, c.Request, userService)
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
	r.DELETE("/community/post", func(c *gin.Context) {
		controllers.DeletePost(c.Writer, c.Request, communityService)
	})
	r.POST("/community/like", func(c *gin.Context) {
		controllers.LikePost(c.Writer, c.Request, communityService)
	})
	r.DELETE("/community/unlike", func(c *gin.Context) {
		controllers.CancelLikePost(c.Writer, c.Request, communityService)
	})
	r.GET("/community/post/likes", func(c *gin.Context) {
		controllers.GetLikesCount(c.Writer, c.Request, communityService)
	})

	// 启动 HTTP 服务
	fmt.Println("Server running on", config.ServerAddress)

	// 启动 HTTP 服务并监听端口
	if err := r.Run(config.ServerAddress); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
