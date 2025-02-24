package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware 验证用户是否已登录
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}
