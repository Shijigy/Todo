package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 注册新用户
func Register(c *gin.Context, userService services.UserService, user models.User) {
	err := userService.RegisterUser(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, models.Response{Status: "success", Message: "User registered successfully"})
}

// Login 用户登录
func Login(c *gin.Context, userService services.UserService) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: "Invalid input"})
		return
	}

	// 调用服务层处理登录
	token, err := userService.LoginUser(context.Background(), loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{Error: err.Error()})
		return
	}

	// 返回生成的 token
	c.JSON(http.StatusOK, models.Response{Message: "Login successful", Data: map[string]string{"token": token}})
}
