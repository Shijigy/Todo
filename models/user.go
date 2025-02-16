package models

import (
	"ToDo/dao"
	"fmt"
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
// CreateATodo 创建user
func CreateAUser(user *User) (err error) {
	err = dao.DB.Table("users").Create(&user).Error
	return
}

func GetAllUser() (userList []*User, err error) {
	if err = dao.DB.Table("users").Find(&userList).Error; err != nil {
		return nil, err
	}
	return
}

func FindAUserByName(username string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Table("users").Where("username=?", username).First(user).Error; err != nil {
		return nil, err
	}
	return
}
func FindAUserByEmail(email string) (user *User, err error) {
	user = new(User)
	if err = dao.DB.Debug().Table("users").Where("email=?", email).First(user).Error; err != nil {
		return nil, err
	}
	return
}
func UpdateUserPasswordByEmail(email string, password string) error {
	err := dao.DB.Table("users").Where("email = ?", email).Update("password", password).Error
	user := new(User)
	if err = dao.DB.Debug().Table("users").Where("email=?", email).First(&user).Error; err != nil {

	}
	fmt.Println(666)
	fmt.Println(user.Username)
	if err != nil {
		return err
	}
	return nil
}
func UpdateUserStatus(email string, status int) error {
	err := dao.DB.Table("users").Where("email = ?", email).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
func DeleteATodo(id string) (err error) {
	err = dao.DB.Table("users").Where("id=?", id).Delete(&User{}).Error
	return
}
func GetStatusByEmail(email string) (int, string) {
	user := new(User)
	dao.DB.Debug().Table("users").Where("email=?", email).First(&user)
	//fmt.Println(user.Status)
	return user.Status, ""
}
