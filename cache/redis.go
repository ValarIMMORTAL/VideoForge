package cache

import (
	"fmt"
	"github.com/pule1234/VideoForge/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	conn *redis.Client
}

func NewRedisConn() (*Redis, error) {
	conf, err := config.LoadConfig("../")
	fmt.Println("conf.RedisSource =" + conf.RedisSource)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisSource,
		Password: conf.RedisPassword,
		DB:       0,
	})

	return &Redis{rdb}, nil
}

//todo 分装一套redis请求方法
