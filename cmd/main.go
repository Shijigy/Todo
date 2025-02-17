package main

import (
	"ToDo/config"
	"ToDo/controllers"
	"ToDo/dao"
	"ToDo/repositories"
	"ToDo/services"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var db *gorm.DB

var checkinRepo repositories.CheckinRepository
var checkinService services.CheckinService
var communityRepo repositories.CommunityRepository
var communityService services.CommunityService
var todoRepo repositories.TodoRepository
var todoService services.TodoService
var likeRepo repositories.LikeRepository

// 启动入口
func main() {
	// 加载配置
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置时出错: %v", err)
	}

	// 输出加载的配置
	fmt.Println("服务器地址:", config.ServerAddress)
	fmt.Println("数据库主机:", config.Database.Host)

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
		log.Fatalf("连接数据库时出错: %v", err)
	}
	defer db.Close()

	// 初始化各个仓库和服务
	checkinRepo = repositories.NewCheckinRepository(db)
	checkinService = services.NewCheckinService(checkinRepo)

	communityRepo = repositories.NewCommunityRepository(db)
	likeRepo = repositories.NewLikeRepository(db)
	communityService = services.NewCommunityService(communityRepo, likeRepo)

	todoRepo = repositories.NewTodoRepository(db)
	todoService = services.NewTodoService(todoRepo)

	// 初始化 Gin 路由器
	r := gin.Default()
	// 设置 session 存储
	store := cookie.NewStore([]byte("secret-key"))
	store.Options(sessions.Options{
		MaxAge:   3600, // 设置过期时间为1小时
		HttpOnly: true, // 设置仅 HTTP 访问
	})
	// 注册会话中间件
	r.Use(sessions.Sessions("session", store), func(c *gin.Context) {
		// 获取当前会话
		session := sessions.Default(c)

		// 如果会话已过期，则重新设置会话信息
		if session.Get("username") == nil {
			// 设置会话值
			session.Set("username", "exampleuser")

			// 设置会话过期时间为1小时
			session.Options(sessions.Options{
				MaxAge:   3600,
				HttpOnly: true,
			})

			// 保存会话
			err := session.Save()
			if err != nil {
				// 处理保存会话时的错误
				c.String(http.StatusInternalServerError, "Failed to save session")
				return
			}
		}
	})

	//创建连接数据库
	err = dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close() // 程序退出关闭数据库连接

	r.POST("/register", controllers.UserRegister)
	// 用户登录接口
	r.POST("/login", controllers.UserLogin)
	// 发送注册验证码
	r.POST("/register-email", controllers.SendEmailRegister)
	// 发送重置密码验证码
	r.POST("/reset-email", controllers.SendEmailReSet)
	// 验证身份接口
	r.POST("/VerifyCode-email", controllers.ResetCodeVerify)
	// 重设密码接口
	r.POST("/reset-password", controllers.ResetPassword)
	// 注销账号
	r.POST("/deactivate", controllers.DeactivateAccount)

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
	fmt.Println("服务器正在运行在", ":8080")

	// 启动 HTTP 服务并监听端口
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器时出错: %v", err)
	}
}
