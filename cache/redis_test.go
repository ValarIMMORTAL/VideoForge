package cache

import (
	"context"
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	InitRedis()

	// 添加元素
	RedisConn.SAdd(ctx, "myset", "value1", "value2", "value3")

	// 查看所有元素
	values, err := RedisConn.SMembers(ctx, "myset").Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Available values:", values)

	// 指定使用某个元素（例如 value2）
	RedisConn.SRem(ctx, "myset", "value2")

	// 再次查看所有元素
	remainingValues, err := RedisConn.SMembers(ctx, "myset").Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Remaining values:", remainingValues)
}
