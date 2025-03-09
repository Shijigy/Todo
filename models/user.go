package models

import (
	"ToDo/dao"
	"fmt"
	"github.com/jinzhu/gorm"
	"image"
	"mime/multipart"
	"time"
)

// User 用户模型
type User struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	Email     string    `json:"email" bson:"email"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Status    int       `json:"status"`
	AvatarURL string    `json:"avatar_url" bson:"avatar_url"` // 头像 URL
}

// CreateAUser 创建user
func CreateAUser(user *User) (err error) {
	err = dao.DB.Table("users").Create(&user).Error
	return
}

// FindAUserByName 根据用户名查询用户
func FindAUserByName(username string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Table("users").Where("username=?", username).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, fmt.Errorf("user not found with username: %s", username)
		}
		return nil, err
	}
	return
}

// FindAUserByEmail 根据邮箱查询用户
func FindAUserByEmail(email string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Table("users").Where("email=?", email).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, err
	}
	return
}

// FindAUserByID 根据用户ID查询用户
func FindAUserByID(userID string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Table("users").Where("id=?", userID).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, fmt.Errorf("user not found with ID: %s", userID)
		}
		return nil, err
	}
	return
}

// UpdateUserPasswordByEmail 根据邮箱更新用户密码
func UpdateUserPasswordByEmail(email string, password string) error {
	err := dao.DB.Table("users").Where("email = ?", email).Update("password", password).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserStatus 根据邮箱更新用户状态
func UpdateUserStatus(email string, status int) error {
	err := dao.DB.Table("users").Where("email = ?", email).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteATodo 根据ID删除用户
func DeleteATodo(id string) (err error) {
	err = dao.DB.Table("users").Where("id=?", id).Delete(&User{}).Error
	return
}

// GetStatusByEmail 根据邮箱获取用户状态
func GetStatusByEmail(email string) (int, string) {
	user := new(User)
	err := dao.DB.Debug().Table("users").Where("email=?", email).First(&user).Error
	if err != nil {
		return 0, fmt.Sprintf("Error fetching user status: %v", err)
	}
	return user.Status, ""
}

// DeleteAUser 删除用户及其所有相关数据
func DeleteAUser(userID string) error {
	// 删除用户记录
	err := dao.DB.Table("users").Where("id = ?", userID).Delete(&User{}).Error
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}

	return nil
}

// isValidImage 验证图片类型和大小
func isValidImage(file *multipart.FileHeader) bool {
	// 打开文件读取
	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	// 检查文件类型
	img, _, err := image.Decode(f)
	if err != nil {
		return false
	}

	// 可选：限制图片的最大尺寸，例如宽度不超过 500px
	maxWidth := 500
	if img.Bounds().Max.X > maxWidth {
		return false
	}

	return true
}

// UpdateUserProfile 更新用户的用户名和头像
func UpdateUserProfile(user *User) error {
	// 更新数据库中的用户名和头像
	err := dao.DB.Table("users").Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username":   user.Username,
		"avatar_url": user.AvatarURL,
		"updated_at": time.Now(),
	}).Error

	if err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}
	return nil
}

// DeleteUserLikes 删除用户的点赞记录
func DeleteUserLikes(userID string) error {
	err := dao.DB.Table("likes").Where("user_id = ?", userID).Delete(&Like{}).Error
	if err != nil {
		return fmt.Errorf("删除点赞记录失败: %v", err)
	}
	return nil
}

// DeleteUserCommunityPosts 删除用户的社区动态
func DeleteUserCommunityPosts(userID string) error {
	err := dao.DB.Table("community_posts").Where("user_id = ?", userID).Delete(&CommunityPost{}).Error
	if err != nil {
		return fmt.Errorf("删除社区动态数据失败: %v", err)
	}
	return nil
}

// DeleteUserCheckins 删除用户的打卡记录
func DeleteUserCheckins(userID string) error {
	err := dao.DB.Table("checkins").Where("user_id = ?", userID).Delete(&Checkin{}).Error
	if err != nil {
		return fmt.Errorf("删除打卡数据失败: %v", err)
	}
	return nil
}

// DeleteUserTodos 删除用户的待办任务记录
func DeleteUserTodos(userID string) error {
	err := dao.DB.Table("todos").Where("user_id = ?", userID).Delete(&Todo{}).Error
	if err != nil {
		return fmt.Errorf("删除待办任务数据失败: %v", err)
	}
	return nil
}

// DeleteUserComments 删除用户的评论记录
func DeleteUserComments(userID string) error {
	// 删除该用户的所有评论记录
	err := dao.DB.Table("comments").Where("user_id = ?", userID).Delete(&Comment{}).Error
	if err != nil {
		return fmt.Errorf("删除评论记录失败: %v", err)
	}
	return nil
}
