package controllers

import (
	middles "ToDo/middlewares"
	"ToDo/models"
	"ToDo/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		restBeanRegister = models.SuccessRestBeanWithData("注册成功")
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

	var restBeanLogin *models.RestBean
	result := services.UserLogin(requestData.Email, requestData.Password)
	if result == "" {
		restBeanLogin = models.SuccessRestBeanWithData("登录成功")
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
	fmt.Println("======================================================")
	fmt.Println(requestData)
	fmt.Println("======================================================")
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
	//sessionID := middles.GetSessionId(c)
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

	//sessionID, _ := c.Cookie("session")
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
	//var emailinit  =  strings.Index(newsession,"reset-password")
	//for i:=emailinit+14;i<len(sessionID)
	var restBeanRegister *models.RestBean
	//if email == nil {
	//	restBeanRegister = models.FailureRestBeanWithData(http.StatusBadRequest, "清先验证邮箱身份")
	//status, _ := models.GetStatusByEmail(emailString)
	//fmt.Println(emailString)
	//fmt.Println(status)
	if emailString != "" {
		services.ResetPassword(requestData.Password, emailString)
		if services.ResetPassword(requestData.Password, emailString) == "" {
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
