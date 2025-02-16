package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetOnesessionNew() {
	router := gin.Default()

	// 创建会话存储
	store := cookie.NewStore([]byte("secret-key"))
	store.Options(sessions.Options{
		MaxAge:   3600, // 设置过期时间为1小时，单位为秒
		HttpOnly: true, // 设置仅HTTP访问
	})

	// 注册会话中间件
	router.Use(sessions.Sessions("session", store))

	router.GET("/set-session", func(c *gin.Context) {
		// 获取当前会话
		session := sessions.Default(c)

		// 设置会话值
		session.Set("username", "exampleuser")

		// 保存会话
		err := session.Save()
		if err != nil {
			// 处理保存会话时的错误
			c.String(http.StatusInternalServerError, "Failed to save session")
			return
		}

		fmt.Println(session.Get("username"))
		c.String(http.StatusOK, "4555")
	})

}
