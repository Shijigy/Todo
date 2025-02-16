package dao

//import "github.com/go-redis/redis"
//
//func RedisCilent(*redis.Client, error) {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379", // Redis 服务器地址
//		Password: "",               // Redis 访问密码（如果有的话）
//		DB:       0,                // Redis 数据库索引
//	})
//
//	// 测试连接是否成功
//	_, err := client.Ping().Result()
//	if err != nil {
//		return nil, err
//	}
//
//	return client, nil
//}
import (
	"github.com/go-redis/redis"
)

func ConnectToRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // Redis 访问密码（如果有的话）
		DB:       0,                // Redis 数据库索引
	})

	// 测试连接是否成功
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
