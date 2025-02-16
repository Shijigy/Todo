package middlewares

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Email struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"email"`
}

// 发送邮箱验证码
func SendCode(email string) string {
	// 读取配置文件
	configData, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("无法读取配置文件:", err)
	}

	// 解析配置文件
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal("无法解析配置文件:", err)
	}
	//发送对象
	recipient := email
	// 生成验证码
	verificationCode := generateVerificationCode()

	// 构建邮件内容
	subject := "验证码"
	body := fmt.Sprintf("你的验证码是：%s", verificationCode)

	// 创建邮件消息
	message := gomail.NewMessage()
	message.SetHeader("From", config.Email.Username)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	// 创建SMTP客户端
	dialer := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password)

	// 发送邮件
	err = dialer.DialAndSend(message)
	if err != nil {
		fmt.Println("发送邮件失败:", err)
		return ""
	}

	return verificationCode
}

func generateVerificationCode() string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成6位验证码
	code := rand.Intn(899999) + 100000

	// 将验证码转换为字符串
	codeStr := strconv.Itoa(code)

	return codeStr

}
