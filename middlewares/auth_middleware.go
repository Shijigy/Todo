package middlewares

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
	var request struct {
		User        models.User `json:"user"`
		CaptchaCode string      `json:"captcha"`
	}

	// 解析请求体
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// 校验邮箱是否已被注册
	existingUser, err := userService.GetUserByEmail(context.Background(), request.User.Email)
	if err == nil && existingUser != nil {
		respondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	// 生成并发送验证码到用户邮箱
	captchaCode, err := utils.GenerateCaptcha()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating captcha")
		return
	}
	err = emailService.SendCaptcha(request.User.Email, captchaCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error sending captcha")
		return
	}

	// 验证验证码是否正确
	if request.CaptchaCode != captchaCode {
		respondWithError(w, http.StatusForbidden, "Invalid captcha")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(request.User.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	request.User.Password = hashedPassword

	// 调用用户服务注册用户
	err = userService.RegisterUser(context.Background(), request.User)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateJWTToken(&request.User)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating JWT")
		return
	}

	// 返回注册成功的信息和 JWT Token
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"token":   token,
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
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	// 调用服务层逻辑处理登录
	ctx := context.Background()
	token, err := userService.LoginUser(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 返回生成的 JWT Token
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

// 封装的错误响应函数
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.Response{Error: message})
}

// 封装的成功响应函数
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
