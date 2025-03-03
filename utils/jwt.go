package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("server")

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成 JWT
func GenerateToken(id string, username string) (accessToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // 令 token 24 小时有效
	claims := &Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	// 加密并获得完整的编码后的字符串token
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return accessToken, err
}

// ParseToken 解析和验证 JWT
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// ParseRefreshToken 验证 JWT
func ParseRefreshToken(aToken string) (newAToken string, err error) {
	accessClaim, err := ParseToken(aToken)
	if err != nil {
		return
	}

	if accessClaim.ExpiresAt > time.Now().Unix() {
		// 如果 access_token 没过期,每一次请求都刷新 refresh_token 和 access_token
		return GenerateToken(accessClaim.ID, accessClaim.Username)
	}

	// 如果过期了,重新登陆
	return "", errors.New("身份过期，重新登陆")
}
