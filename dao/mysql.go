package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	DB *gorm.DB
)

// Config 用来存储从配置文件读取的数据库配置
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Charset  string `yaml:"charset"`
	} `yaml:"database"`
}

var config Config

// LoadConfig 读取配置文件
func LoadConfig() error {
	// 读取配置文件
	data, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		return fmt.Errorf("无法读取配置文件: %v", err)
	}
	// 解析配置文件
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("无法解析配置文件: %v", err)
	}
	return nil
}

// InitMySQL 初始化数据库连接
func InitMySQL() (err error) {
	// 加载配置
	if err := LoadConfig(); err != nil {
		return err
	}

	// 拼接数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Dbname,
		config.Database.Charset,
	)

	// 连接数据库
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
