package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// GenerateSHA256Hash 生成 SHA-256 哈希
func GenerateSHA256Hash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// GenerateSHA512Hash 生成 SHA-512 哈希
func GenerateSHA512Hash(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// HashPassword 使用 bcrypt 对密码进行加密
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswordHash 比较明文密码和加密后的密码是否一致
func ComparePasswordHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
