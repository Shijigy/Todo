package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

//	func GetSessionId(w http.ResponseWriter, r *http.Request) string {
//		session, err := r.Cookie("")
//		if err != nil {
//			// 处理获取会话时的错误
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return "错误"
//		}
//
//		// 从会话 cookie 中获取会话标识符
//		sessionID := session.Value
//		return sessionID
//
// }

func GetSessionId(c *gin.Context) string {
	SetSessionId(c)
	session, err := c.Cookie("session")
	if err != nil {
		// 处理获取会话时的错误
		log.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)

	}

	// 从会话 cookie 中获取会话标识符
	sessionID := session
	//c.String(http.StatusOK, sessionID)
	return sessionID
}
func SetSessionId(c *gin.Context) {
	// 生成一个具有足够熵的随机字节数组
	bytes := make([]byte, 32)

	// 将随机字节数组转换为 base64 编码字符串
	_, err := rand.Read(bytes)
	if err != nil {
		// 处理生成随机字节数组时的错误
		log.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	sessionID := base64.URLEncoding.EncodeToString(bytes)
	// 生成会话标识符
	//sessionID := "your-session-id"

	// 设置会话cookie
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(c.Writer, cookie)

	//c.String(http.StatusOK, "会话cookie已设置")
}

func SetSessionAttribute(c *gin.Context, key string, value interface{}) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

func GetSessionAttribute(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	value := session.Get(key)
	return value
}
func DeleteSessionKey(c *gin.Context, key string) {
	// 获取会话对象
	session := sessions.Default(c)

	// 移除会话属性
	session.Delete(key)
	session.Save()

}
func SetSessionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取或设置 session ID
		sessionID := GetSessionId(c)

		// 将 session ID 设置到请求上下文中
		c.Set("session", sessionID)

		// 调用下一个处理程序
		c.Next()
	}
}
