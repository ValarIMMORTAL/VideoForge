package cache

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	RedisClient *Redis
	RedisConn   *redis.Client // 全局使用rabbitmq连接
	once        sync.Once     // 确保只初始化一次
)

// main中调用此函数，初始化全局rabbitmq变量
func InitRedis() error {
	var err error
	once.Do(func() {
		RedisClient, err = NewRedisConn()
		RedisConn = RedisClient.conn
	})

	return err
}

// GetRabbitMQ 获取全局 RabbitMQ 实例
func GetRedisConn() *Redis {
	return RedisClient
}
