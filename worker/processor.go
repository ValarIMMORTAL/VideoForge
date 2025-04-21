package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor struct {
	mux      *asynq.ServeMux
	registry map[string]func(context.Context, interface{}) error
}

func NewTaskProcessor() *TaskProcessor {
	mux := asynq.NewServeMux()
	processor := &TaskProcessor{
		mux:      mux,
		registry: make(map[string]func(context.Context, interface{}) error),
	}
	//动态任务
	mux.HandleFunc("dynamic_task", processor.handleDynamicTask)
	return processor
}

// 注册任务
func (p *TaskProcessor) Register(funcName string, handler func(context.Context, interface{}) error) {
	p.registry[funcName] = handler
}

func (p *TaskProcessor) handleDynamicTask(ctx context.Context, task *asynq.Task) error {
	var dynTask DynamicTask
	if err := json.Unmarshal(task.Payload(), &dynTask); err != nil {
		return fmt.Errorf("failed to unmarshal dynamic task: %w", err)
	}

	handler, ok := p.registry[dynTask.FuncName]
	if !ok {
		return fmt.Errorf("no handler registered for function: %s", dynTask.FuncName)
	}

	return handler(ctx, dynTask.Payload)
}

func (p *TaskProcessor) Start(redisOpt asynq.RedisClientOpt) error {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})
	// 运行任务
	return server.Run(p.mux)
}
