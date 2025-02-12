package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware(c *gin.Context) {
	// 记录请求信息
	start := time.Now()
	log.Printf("Started %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())

	// 调用下一个处理器
	c.Next()

	// 记录处理完请求后的信息
	duration := time.Since(start)
	log.Printf("Completed %s %s in %v", c.Request.Method, c.Request.URL.Path, duration)
}
