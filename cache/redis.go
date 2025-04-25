package cache

import (
	"context"
	"fmt"
	"github.com/pule1234/VideoForge/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Conn *redis.Client
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

// 向集合中添加数据
func (r *Redis) SAdd(ctx context.Context, key string, members ...interface{}) {
	r.Conn.SAdd(ctx, key, members)
}

// 查看所有元素
func (r *Redis) SMembers(ctx context.Context, key string) {
	r.Conn.SMembers(ctx, key)
}

// 删除指定的数据
func (r *Redis) SRem(ctx context.Context, key string, members ...interface{}) {
	r.Conn.SRem(ctx, key, members)
}

// hash
func (r *Redis) HSet(ctx context.Context, key, field string, value interface{}) {
	r.Conn.HSet(ctx, key, field, value)
}

func (r *Redis) Scan(ctx context.Context, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error) {
	keys, newCursor, err = r.Conn.Scan(ctx, cursor, "retry:*", 100).Result()
	return
}

func (r *Redis) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (msgs []string, err error) {
	msgs, err = r.Conn.ZRangeByScore(ctx, key, opt).Result()
	return
}

func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.Conn.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
