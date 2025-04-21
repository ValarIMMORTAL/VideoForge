package worker

import "github.com/hibiken/asynq"

type TaskDistributor interface {
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisopt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisopt)
	return &RedisTaskDistributor{
		client: client,
	}
}
