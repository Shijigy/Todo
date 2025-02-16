package services

import (
	"ToDo/dao"
	middles "ToDo/middlewares"
	"ToDo/models"
	"fmt"
	"strconv"
	"time"
)

// var PasswordMap = make(map[string]string])

// UserRegister 普通用户的注册
func UserRegister(username string, password string, email string, code string, sessionId string) string {
	//连接redis
	RedisClient, err := dao.ConnectToRedis()
	// 检查Redis中键是否存在
	key := "email:" + sessionId + ":" + email + ":false"
	fmt.Println(key)
	exist, err := RedisClient.Exists(key).Result()
	if err != nil {
		return "内部错误"
	}
	if exist == 0 {
		return "请先请求一封验证码邮件"
	}
	// 获取Redis中键对应的值
	result, err := RedisClient.Get(key).Result()
	if result == "" {
		return "验证码失效，请重新请求"
	}
	if result == code {
		users, _ := models.FindAUserByName(username)
		if users != nil {
			return "此用户名已被注册，请更换用户名"
		}
		RedisClient.Del(key)
		// 对密码进行加密
		//privateencode := middles.Encode(password)
		privateencode := password
		password = privateencode
		// 创建新用户
		user := models.User{
			Username: username,
			Password: password,
			Email:    email,
			Status:   0,
		}

		err := models.CreateAUser(&user)
		if err != nil {
			return "内部错误"
		}

		return "" // 注册成功，返回空字符串表示成功
	} else {
		return "验证码错误，请检查后再提交"
	}
}

// UserLogin 普通用户的登录
func UserLogin(email string, password string) string {
	//判断用户名是不是为空
	if email == "" {
		return "邮箱不能为空"
	}
	// 根据用户名从数据库中获取用户信息
	user, _ := models.FindAUserByEmail(email)

	// 验证用户是否存在
	if user == nil {
		return "用户不存在"
	}
	if password != user.Password {
		return "密码不正确"
	}

	return "" // 登录成功，返回空字符串表示成功
}

// SendEmail 发送验证码
func SendEmail(email string, sessionId string, hashAccount bool) string {
	//连接redis
	RedisClient, err := dao.ConnectToRedis()
	key := "email:" + sessionId + ":" + email + ":" + strconv.FormatBool(hashAccount)
	pan, _ := RedisClient.Exists(key).Result()
	if pan == 1 {
		expire, _ := RedisClient.TTL(key).Result()
		if expire > 120*time.Second {
			return "请求频繁，请稍后再试"
		}
	}

	// 模拟查找账户
	account, _ := models.FindAUserByEmail(email)
	if hashAccount && account == nil {
		return "没有此邮件地址的账户"
	}
	if !hashAccount && account != nil {
		return "此邮箱已被其他用户注册"
	}

	// 模拟发送邮件
	result := middles.SendCode(email)
	if result == "" {
		return "邮件发送失败，请检查邮件地址是否有效"
	}

	err = RedisClient.Set(key, result, 3*time.Minute).Err()
	if err != nil {
		// 处理缓存错误
	}

	return ""
}

// ResetCode 修改时申请发送验证码，就是验证邮箱那一步，验证完后才能开始修改密码
func ResetCode(email string, code string, sessionId string) string {
	//连接redis
	RedisClient, _ := dao.ConnectToRedis()
	// 检查Redis中键是否存在
	key := "email:" + sessionId + ":" + email + ":true"
	exist, _ := RedisClient.Exists(key).Result()
	if exist == 0 {
		return "请先请求一封验证码邮件"
	}

	// 获取存储在Redis中的值
	value, err := RedisClient.Get(key).Result()
	if err != nil {
		// 处理错误
		return "获取Redis中的值时出错"
	}

	if value == "" {
		return "验证码失效，请重新请求"
	}

	if value == code {
		//设置改密码权限
		models.UpdateUserStatus(email, 1)
		// 删除Redis中的键
		_, err := RedisClient.Del(key).Result()
		if err != nil {
			// 处理错误
			return "删除Redis中的键时出错"
		}

		return "" // 返回空表示验证通过
	} else {
		return "验证码错误，请检查后再提交"
	}
}

// ResetPassword 邮箱验证通过后才能修改密码
func ResetPassword(email string, password string) string {
	//password = middles.Encode(password)
	models.UpdateUserPasswordByEmail(password, email)
	fmt.Println(password)
	models.UpdateUserStatus(email, 0)
	return ""
}
