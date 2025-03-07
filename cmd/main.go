package main

import (
	"ToDo/controllers"
	"ToDo/dao"
	"ToDo/repositories"
	"ToDo/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var checkinRepo repositories.CheckinRepository
var checkinService services.CheckinService
var communityRepo repositories.CommunityRepository
var communityService services.CommunityService
var todoRepo repositories.TodoRepository
var todoService services.TodoService
var likeRepo repositories.LikeRepository

// 启动入口
func main() {

	//创建连接数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close() // 程序退出关闭数据库连接

	// 初始化各个仓库和服务
	checkinRepo = repositories.NewCheckinRepository(dao.DB)
	checkinService = services.NewCheckinService(checkinRepo)

	communityRepo = repositories.NewCommunityRepository(dao.DB)
	likeRepo = repositories.NewLikeRepository(dao.DB)
	communityService = services.NewCommunityService(communityRepo, likeRepo)

	todoRepo = repositories.NewTodoRepository(dao.DB)
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
	// 更新用户名和头像
	r.PUT("/profile", controllers.UpdateUserInfoController)

	// Todo 路由
	// 创建任务
	r.POST("/create-todo", func(c *gin.Context) {
		controllers.CreateTodo(c.Writer, c.Request, todoService)
	})
	// 获取所有任务
	r.GET("/get-todo", func(c *gin.Context) {
		controllers.GetTodos(c.Writer, c.Request, todoService)
	})

	// 修改任务（通过请求体传递任务 ID 和其他字段）
	r.PUT("/reset-todo", func(c *gin.Context) {
		controllers.UpdateTodo(c.Writer, c.Request, todoService)
	})
	// 删除任务（通过请求体传递任务 ID）
	r.DELETE("/delete-todo/:id", func(c *gin.Context) {
		controllers.DeleteTodo(c, todoService)
	})
	// 标记任务为已完成
	r.PUT("/complete", func(c *gin.Context) {
		controllers.MarkTodoAsCompleted(c, todoService)
	})

	// Checkin 路由
	r.POST("/checkin", func(c *gin.Context) {
		controllers.Checkin(c.Writer, c.Request, checkinService)
	})
	r.GET("/get-checkin", func(c *gin.Context) {
		controllers.GetCheckinRecordByUserID(c.Writer, c.Request, checkinService)
	})
	r.PUT("/checkin/complete", func(c *gin.Context) {
		controllers.MarkCheckinComplete(c.Writer, c.Request, checkinService)
	})
	r.PUT("/checkin/update-count", func(c *gin.Context) {
		controllers.UpdateCheckinCount(c.Writer, c.Request, checkinService)
	})
	r.DELETE("/checkin/delete", func(c *gin.Context) {
		controllers.DeleteCheckin(c.Writer, c.Request, checkinService)
	})

	// 社区动态路由
	r.POST("/community/post", func(c *gin.Context) {
		controllers.CreatePost(c.Writer, c.Request, communityService)
	})
	r.GET("/community/list", func(c *gin.Context) {
		controllers.GetPosts(c.Writer, c.Request, communityService)
	})
	r.GET("/community/all-list", func(c *gin.Context) {
		controllers.GetAllPosts(c.Writer, c.Request, communityService)
	})
	r.DELETE("/community/del-post", func(c *gin.Context) {
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
	r.POST("/comment", func(c *gin.Context) {
		controllers.CreateComment(c.Writer, c.Request, communityService)
	})
	r.DELETE("/del-comments", func(c *gin.Context) {
		controllers.DeleteComment(c.Writer, c.Request, communityService)
	})
	r.GET("/comments", func(c *gin.Context) {
		controllers.GetComments(c.Writer, c.Request, communityService)
	})

	// 启动 HTTP 服务并监听端口
	if err := r.Run("0.0.0.0:9999"); err != nil {
		log.Fatalf("启动服务器时出错: %v", err)
	}
}
