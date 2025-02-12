package controllers

import (
	"ToDo/models"
	"ToDo/services"
	"ToDo/utils"
	"context"
	"encoding/json"
	"net/http"
)

// Register 注册新用户
func Register(w http.ResponseWriter, r *http.Request, userService services.UserService, emailService services.EmailService) {
	var user models.User
	var captchaInput string

	// 解析请求体
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 从请求中获取验证码
	if err := json.NewDecoder(r.Body).Decode(&captchaInput); err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid captcha"})
		return
	}

	// 校验邮箱是否已被注册
	existingUser, err := userService.GetUserByEmail(context.Background(), user.Email)
	if err == nil && existingUser != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.Response{Error: "Email already registered"})
		return
	}

	// 生成并发送验证码到用户邮箱
	captchaCode, err := utils.GenerateCaptcha()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error generating captcha"})
		return
	}
	err = emailService.SendCaptcha(user.Email, captchaCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error sending captcha"})
		return
	}

	// 验证验证码是否正确
	if captchaInput != captchaCode {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid captcha"})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	// 调用用户服务注册用户
	err = userService.RegisterUser(context.Background(), user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateJWTToken(&user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error generating JWT"})
		return
	}

	// 返回注册成功的信息和 JWT Token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Message: "User registered successfully",
		Token:   token,
	})
}

// Login 用户登录
func Login(w http.ResponseWriter, r *http.Request, userService services.UserService) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 解析请求体
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 调用服务层逻辑处理登录
	ctx := context.Background()
	token, err := userService.LoginUser(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		// 返回错误信息，格式化为 JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 返回生成的 JWT Token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Message: "Login successful", Data: map[string]string{"token": token}})
}
