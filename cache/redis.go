package cache

import (
	"context"
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

// todo 分装一套redis请求方法
// 向集合中添加数据
func (r *Redis) SAdd(ctx context.Context, key string, members ...interface{}) {
	r.conn.SAdd(ctx, key, members)
}

// 查看所有元素
func (r *Redis) SMembers(ctx context.Context, key string) {
	r.conn.SMembers(ctx, key)
}

// 删除指定的数据
func (r *Redis) SRem(ctx context.Context, key string, members ...interface{}) {
	r.conn.SRem(ctx, key, members)
}
