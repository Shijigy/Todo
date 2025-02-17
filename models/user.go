package models

import (
	"ToDo/dao"
	"fmt"
	"github.com/jinzhu/gorm"
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
}

/*
	User这个Model的增删改查操作都放在这里
*/

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

	// 删除用户的相关数据，如 Todo、Checkin 等
	// 例如：
	// err = dao.DB.Table("todos").Where("user_id = ?", userID).Delete(&Todo{}).Error
	// if err != nil {
	//     return fmt.Errorf("删除用户的 Todo 数据失败: %v", err)
	// }

	return nil
}
