package config

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"github.com/joho/godotenv"         // 用于加载 .env 文件
	"log"
	"os"
)

// AppConfig 配置结构体
type AppConfig struct {
	ServerAddress string
	Database      DatabaseConfig
	JWTSecretKey  string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
}

// LoadConfig 从环境变量或配置文件加载配置
func LoadConfig() (*AppConfig, error) {
	// 加载 .env 配置文件
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	// 配置结构体
	config := &AppConfig{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),     // MySQL host
			Port:     getEnv("DB_PORT", "3306"),          // MySQL port
			Username: getEnv("DB_USER", "root"),          // MySQL username
			Password: getEnv("DB_PASSWORD", "SHANG6688"), // MySQL password
			DbName:   getEnv("DB_NAME", "todo"),          // MySQL database name
		},
		JWTSecretKey: getEnv("JWT_SECRET_KEY", "defaultsecretkey"),
	}

	return config, nil
}

// getEnv 获取环境变量值，如果没有设置则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// 输出加载的配置，或者将其传递给应用程序
	fmt.Println("Server address:", config.ServerAddress)
	fmt.Println("Database host:", config.Database.Host)
}
