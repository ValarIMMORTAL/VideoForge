package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

type TaskDistributor struct {
	client *asynq.Client
}

type DynamicTask struct {
	FuncName string      `json:"func_name"`
	Payload  interface{} `json:"payload"`
}

func NewRedisTaskDistributor(redisopt asynq.RedisClientOpt) *TaskDistributor {
	client := asynq.NewClient(redisopt)
	return &TaskDistributor{
		client: client,
	}
}

func (d *TaskDistributor) EnqueueDynamicTask(
	ctx context.Context,
	funcName string,
	payload interface{},
	opts ...asynq.Option,
) error {
	dynTask := DynamicTask{ //任务执行的函数名称 及 逻辑
		FuncName: funcName,
		Payload:  payload,
	}

	data, err := json.Marshal(dynTask)
	if err != nil {
		return fmt.Errorf("failed to marshal dynamic task: %w", err)
	}
	task := asynq.NewTask("dynamic_task", data, opts...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	fmt.Printf("Enqueued task: %+v\n", info)
	return nil
}
