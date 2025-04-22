package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

// TaskDistributor 封装任务发布逻辑
type TaskDistributor struct {
	client *asynq.Client
}

// NewTaskDistributor 初始化任务发布器
func NewTaskDistributor(redisOpt asynq.RedisClientOpt) *TaskDistributor {
	return &TaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}

// EnqueueDynamicTask 发送一个通用任务
func (d *TaskDistributor) EnqueueDynamicTask(ctx context.Context, taskName string, payload interface{}, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化 payload 失败: %w", err)
	}

	task := asynq.NewTask(taskName, data, opts...)
	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("任务派发失败: %w", err)
	}

	return nil
}
