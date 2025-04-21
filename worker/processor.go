package worker

import (
	"github.com/hibiken/asynq"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/redis/go-redis/v9"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	// todo  统一的任务处理函数
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) *RedisTaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
			Logger: logger,
		})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
