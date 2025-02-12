package main

import (
	"ToDo/config"
	"ToDo/controllers"
	"ToDo/middlewares"
	"ToDo/repositories"
	"ToDo/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
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

	r := mux.NewRouter()

	r.Use(middlewares.LoggingMiddleware)

	// CORS 中间件，允许前端访问后端 API
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                             // 允许所有的域名进行访问
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // 允许的请求方法
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // 允许的请求头
	})

	// 路由设置
	r.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.Login(w, r, userService)
	}).Methods("POST")
	r.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		emailService := services.NewEmailService()
		controllers.Register(w, r, userService, emailService)
	}).Methods("POST")

	// Todo 路由
	r.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateTodo(w, r, todoService)
	}).Methods("POST")
	r.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetTodo(w, r, todoService)
	}).Methods("GET")
	r.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateTodo(w, r, todoService)
	}).Methods("PUT")
	r.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteTodo(w, r, todoService)
	}).Methods("DELETE")

	// Checkin 路由
	r.HandleFunc("/checkin", func(w http.ResponseWriter, r *http.Request) {
		controllers.Checkin(w, r, checkinService)
	}).Methods("POST")

	// 社区动态路由
	r.HandleFunc("/community/post", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreatePost(w, r, communityService)
	}).Methods("POST")
	r.HandleFunc("/community/list", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetPosts(w, r, communityService)
	}).Methods("GET")
	r.HandleFunc("/community/{post_id}/like", func(w http.ResponseWriter, r *http.Request) {
		// 提取 path 参数 (post_id)
		vars := mux.Vars(r)
		postID := vars["post_id"]

		// 解析请求体中的 UserID
		var requestBody struct {
			UserID string `json:"user_id"` // 修改字段为小写 user_id
		}

		// 解析请求体
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 校验 UserID 是否有效
		if requestBody.UserID == "" {
			http.Error(w, "UserID is required", http.StatusBadRequest)
			return
		}

		// 调用 LikePost 控制器并传入 postID 和 UserID
		controllers.LikePost(w, r, communityService, postID, requestBody.UserID)
	}).Methods("POST")

	// 启动 HTTP 服务
	fmt.Println("Server running on", config.ServerAddress)

	// 使用 CORS 中间件
	handler := corsHandler.Handler(r)

	// 启动 HTTP 服务并监听端口
	log.Fatal(http.ListenAndServe(config.ServerAddress, handler))
}
