package services

import (
	"ToDo/dao"
	"ToDo/models"
	"ToDo/utils"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"net/http"
)

// UserRegister 用户的注册
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
		privateencode := password
		password = privateencode
		// 检查头像URL，如果没有传入则设置为默认头像
		avatarURL := "http://ssjwo2ece.hn-bkt.clouddn.com/post-images/1741094306314528200_KIb9cH" // 默认头像URL
		// 创建新用户
		user := models.User{
			Username:  username,
			Password:  password,
			Email:     email,
			Status:    0,
			AvatarURL: avatarURL, // 设置默认头像
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
func UserLogin(email string, password string) (string, string, string, string) {
	// 判断邮箱是不是为空
	if email == "" {
		return "", "", "", "邮箱不能为空"
	}

	// 根据邮箱从数据库中获取用户信息
	user, err := models.FindAUserByEmail(email)
	if err != nil {
		return "", "", "", "获取用户信息失败"
	}

	// 验证用户是否存在
	if user == nil {
		return "", "", "", "用户不存在"
	}

	// 使用 bcrypt.CompareHashAndPassword 比对用户输入的密码与数据库中的加密密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 如果密码不匹配，返回错误信息
		return "", "", "", "密码不正确"
	}

	// 登录成功，返回用户名、头像URL和用户ID
	return user.Username, user.AvatarURL, user.ID, ""
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
	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "加密密码时发生错误"
	}

	// 将加密后的密码存入数据库
	err = models.UpdateUserPasswordByEmail(email, string(hash))
	if err != nil {
		return "更新密码时发生错误"
	}

	// 打印加密后的密码
	fmt.Println("Encrypted Password:", string(hash))

	// 重置密码修改状态
	models.UpdateUserStatus(email, 0)
	return "" // 返回空表示密码重置成功
}

// DeactivateUser 停用用户并删除所有相关数据
func DeactivateUser(userID string) error {

	// 删除 likes 表中的相关数据
	err := models.DeleteUserLikes(userID)
	if err != nil {
		return fmt.Errorf("删除用户的点赞数据失败: %v", err)
	}

	// 删除 comments 表中的相关数据
	err = models.DeleteUserComments(userID)
	if err != nil {
		return fmt.Errorf("删除用户的评论数据失败: %v", err)
	}

	// 删除 community_posts 表中的相关数据
	err = models.DeleteUserCommunityPosts(userID)
	if err != nil {
		return fmt.Errorf("删除用户的社区动态数据失败: %v", err)
	}

	// 删除 checkins 表中的相关数据
	err = models.DeleteUserCheckins(userID)
	if err != nil {
		return fmt.Errorf("删除用户的打卡数据失败: %v", err)
	}

	// 删除 todos 表中的相关数据
	err = models.DeleteUserTodos(userID)
	if err != nil {
		return fmt.Errorf("删除用户的待办任务数据失败: %v", err)
	}
	// 从数据库中删除用户数据
	err = models.DeleteAUser(userID)
	if err != nil {
		return fmt.Errorf("删除用户数据失败: %v", err)
	}

	return nil
}

// UpdateUserInfo 更新用户信息，包括用户名和头像
func UpdateUserInfo(userID string, username string, avatarFile multipart.File, request *http.Request) (*models.User, error) {
	// 查找用户
	var user models.User
	if err := dao.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 更新用户名，如果传入了用户名
	if username != "" {
		user.Username = username
	}

	// 处理头像文件上传
	if avatarFile != nil {
		avatarURL, err := utils.UploadImageToQiNiu(request)
		if err != nil {
			return nil, errors.New("failed to upload avatar")
		}
		user.AvatarURL = avatarURL
	}

	// 保存更新后的用户信息
	if err := dao.DB.Save(&user).Error; err != nil {
		return nil, errors.New("failed to save user info")
	}

	return &user, nil
}
