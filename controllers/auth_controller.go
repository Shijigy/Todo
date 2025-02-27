package controllers

import (
	middles "ToDo/middlewares"
	"ToDo/models"
	"ToDo/services"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"image"
	"mime/multipart"
	"net/http"
)

// UserRegister 用户注册并返回信息给客户端
func UserRegister(c *gin.Context) {
	sessionID, _ := c.Cookie("session")
	var requestData struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
		Email    string `form:"email" json:"email"`
		Code     string `form:"code" json:"code"`
	}

	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 调用 UserRegister 函数进行注册
	result := services.UserRegister(requestData.Username, string(hash), requestData.Email, requestData.Code, sessionID)

	var restBeanRegister *models.RestBean
	if result == "" {
		// 获取用户注册后的详细信息，如用户名和头像URL
		username, avatarURL, _ := services.UserLogin(requestData.Email, requestData.Password)
		restBeanRegister = models.SuccessRestBeanWithData(gin.H{
			"message":    "注册成功",
			"username":   username,
			"avatar_url": avatarURL,
		})
	} else {
		restBeanRegister = models.FailureRestBeanWithData(http.StatusBadRequest, result)
	}

	// 返回注册结果给前端
	c.JSON(restBeanRegister.Status, restBeanRegister)
}

// 用户登录并返回信息给客户端
func UserLogin(c *gin.Context) {
	var requestData struct {
		Email    string `form:"email" json:"email"`
		Password string `form:"password" json:"password"`
	}

	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 调用登录函数，并获取用户名和头像URL
	username, avatarURL, result := services.UserLogin(requestData.Email, requestData.Password)

	var restBeanLogin *models.RestBean
	if result == "" {
		restBeanLogin = models.SuccessRestBeanWithData(gin.H{
			"message":    "登录成功",
			"username":   username,
			"avatar_url": avatarURL,
		})
	} else {
		restBeanLogin = models.FailureRestBeanWithData(http.StatusBadRequest, result)
	}

	// 返回登录结果给前端
	c.JSON(restBeanLogin.Status, restBeanLogin)
}

// SendEmailRegister 发送注册验证码
func SendEmailRegister(c *gin.Context) {
	sessionID, _ := c.Cookie("session")
	fmt.Println(sessionID)
	//sessionID := middles.GetSessionId(c)
	var requestData struct {
		Email string `form:"email" json:"email"`
	}
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	// 调用UserRegister函数进行注册
	result := services.SendEmail(requestData.Email, sessionID, false)

	// 根据注册结果返回相应的数据给前端
	// 封装注册结果为RestBean对象
	var restBeanRegister *models.RestBean
	if result == "" {
		restBeanRegister = models.SuccessRestBeanWithData("邮件已发送，请注意查收")

	} else {
		restBeanRegister = models.FailureRestBeanWithData(http.StatusBadRequest, result)
	}
	//返回注册结果给前端
	c.JSON(restBeanRegister.Status, restBeanRegister)
}

// SendEmailReSet 发送重置密码验证码
func SendEmailReSet(c *gin.Context) {
	sessionID, _ := c.Cookie("session")
	//sessionID := middles.GetSessionId(c)
	var requestData struct {
		Email string `form:"email" json:"email"`
	}
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	// 调用UserRegister函数进行注册
	result := services.SendEmail(requestData.Email, sessionID, true)

	// 根据注册结果返回相应的数据给前端
	// 封装注册结果为RestBean对象
	var restBeanRegister *models.RestBean
	if result == "" {
		restBeanRegister = models.SuccessRestBeanWithData("邮件已发送，请注意查收")

	} else {
		restBeanRegister = models.FailureRestBeanWithData(http.StatusBadRequest, result)
	}
	//返回注册结果给前端
	c.JSON(restBeanRegister.Status, restBeanRegister)
	c.JSON(http.StatusOK, gin.H{"message": "重置密码验证码已发送"})
}

// ResetCodeVerify 验证邮箱验证码
func ResetCodeVerify(c *gin.Context) {
	sessionID, _ := c.Cookie("session")
	// 获取验证参数
	var requestData struct {
		Email string `form:"email" json:"email"`
		Code  string `form:"code" json:"code"`
	}
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	// 调用函数进行身份验证
	result := services.ResetCode(requestData.Email, requestData.Code, sessionID)

	// 根据注册结果返回相应的数据给前端
	// 封装注册结果为RestBean对象
	var restBeanRegister *models.RestBean
	if result == "" {
		// 在会话中设置重置密码的相关属性
		middles.SetSessionAttribute(c, "reset-password", requestData.Email)
		//email := middles.GetSessionAttribute(c, "reset-password")
		//fmt.Println(email)
		restBeanRegister = models.SuccessRestBean()

	} else {
		restBeanRegister = models.FailureRestBeanWithData(http.StatusBadRequest, result)
	}
	//返回结果给前端
	c.JSON(restBeanRegister.Status, restBeanRegister)
	c.JSON(http.StatusOK, gin.H{"message": "验证码验证通过"})
}

// ResetPassword 重设密码
func ResetPassword(c *gin.Context) {

	// 获取重置密码参数
	var requestData struct {
		Password string `form:"password" json:"password"`
		//Email    string `form:"email" json:"email"`
	}

	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	//从会话中获取重置密码的属性
	email := middles.GetSessionAttribute(c, "reset-password")
	//将邮箱地址转换为字符串
	emailString, _ := email.(string)

	var restBeanRegister *models.RestBean
	if emailString != "" {
		services.ResetPassword(emailString, requestData.Password)
		if services.ResetPassword(emailString, requestData.Password) == "" {
			middles.DeleteSessionKey(c, "reset-password")
			fmt.Println(requestData.Password)
			restBeanRegister = models.SuccessRestBeanWithData("密码重置成功")
			models.UpdateUserStatus(emailString, 0)
		} else {
			restBeanRegister = models.FailureRestBeanWithData(500, "内部错误")
		}
	} else {
		restBeanRegister = models.FailureRestBeanWithData(http.StatusBadRequest, "请先完成邮箱认证")
	}

	//返回结果给前端
	c.JSON(restBeanRegister.Status, restBeanRegister)
}

// DeactivateAccount 用户注销并删除所有数据
func DeactivateAccount(c *gin.Context) {
	var requestData struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}

	// 获取请求数据
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 如果用户名为空，则返回错误
	if requestData.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	// 从数据库获取该用户的加密密码
	user, err := models.FindAUserByName(requestData.Username) // 根据请求中的用户名查找用户
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息获取失败"})
		return
	}

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码验证失败"})
		return
	}
	if requestData.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码验证失败"})
		return
	}

	// 用户验证通过，删除用户数据
	err = services.DeactivateUser(user.ID) // 通过 user.ID 删除用户数据
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注销账户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "账户已成功注销"})
}

// isValidImage 验证图片类型和大小
func isValidImage(file *multipart.FileHeader) bool {
	// 打开文件读取
	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	// 检查文件类型
	img, _, err := image.Decode(f)
	if err != nil {
		return false
	}

	// 可选：限制图片的最大尺寸，例如宽度不超过 500px
	maxWidth := 500
	if img.Bounds().Max.X > maxWidth {
		return false
	}

	return true
}

// UpdateProfile 处理更新用户名和头像请求
func UpdateProfile(c *gin.Context) {
	// 1. 获取登录用户的ID
	userID, exists := c.Get("user_id") // 假设在用户登录后，用户ID存储在上下文中
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	var requestData struct {
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}

	// 2. 绑定请求数据
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 3. 调用服务层函数来更新用户名和头像
	result := services.UpdateUserProfile(userID.(string), requestData.Username, requestData.AvatarURL)

	if result == "" {
		c.JSON(http.StatusOK, gin.H{"message": "个人信息更新成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": result})
	}
}
