package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

// InitMySQL 初始化数据库连接
func InitMySQL() (err error) {
	// 数据库连接字符串
	dsn := "root:SHANG6688@tcp(127.0.0.1:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("无法连接数据库: %v", err)
	}

	// 测试数据库连接
	if err := DB.DB().Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}
	return nil
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		sqlDB := DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}
