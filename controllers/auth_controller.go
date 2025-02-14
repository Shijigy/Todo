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
	var req struct {
		User         models.User `json:"user"`
		CaptchaInput string      `json:"captchaInput"`
	}

	// 解码请求体
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid input"})
		return
	}

	// 检查邮箱是否已注册
	existingUser, err := userService.GetUserByEmail(context.Background(), req.User.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error checking email"})
		return
	}
	if existingUser != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.Response{Error: "Email already registered"})
		return
	}

	// 生成并发送验证码
	captchaCode, err := utils.GenerateCaptcha()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error generating captcha"})
		return
	}
	err = emailService.SendCaptcha(req.User.Email, captchaCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error sending captcha"})
		return
	}

	// 验证验证码
	if req.CaptchaInput != captchaCode {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(models.Response{Error: "Invalid captcha"})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.User.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error hashing password"})
		return
	}
	req.User.Password = hashedPassword

	// 注册用户
	err = userService.RegisterUser(context.Background(), req.User)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	// 生成 JWT token
	token, err := utils.GenerateJWTToken(&req.User)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Error: "Error generating JWT"})
		return
	}

	// 返回成功信息和 token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
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
