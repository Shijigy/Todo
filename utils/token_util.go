package utils

import (
	"ToDo/models"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtKey = []byte("your_secret_key")

// Claims 用于 JWT 中存储用户信息
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成 JWT Token
func GenerateJWTToken(user *models.User) (string, error) {
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
		},
	}

	// 创建一个 token 实例
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 用密钥签名 token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken 验证 JWT Token 是否有效
func ValidateJWTToken(tokenString string) (string, error) {
	// 解析 token，获取 claims 信息
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// 验证 token 是否有效
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("Invalid token signature")
		}
		return "", fmt.Errorf("Failed to parse token: %v", err)
	}

	// 验证 token 是否有效
	if !token.Valid {
		return "", fmt.Errorf("Token is not valid")
	}

	// 返回用户名
	return claims.Username, nil
}
