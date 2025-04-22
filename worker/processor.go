package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"reflect"
)

type TaskProcessor struct {
	server      *asynq.Server
	mux         *asynq.ServeMux
	taskHandler map[string]func(context.Context, any) error
}

// NewTaskProcessor 初始化处理器
func NewTaskProcessor(redisOpt asynq.RedisClientOpt) *TaskProcessor {
	return &TaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{
			Concurrency: 10,
		}),
		mux:         asynq.NewServeMux(),
		taskHandler: make(map[string]func(context.Context, any) error),
	}
}

// 注册任务处理函数
// handlerFunc : func(context.Context, Struct) error   例如： func(ctx context.Context, payload SendEmailPayload) error
func (p *TaskProcessor) Register(taskName string, handlerFunc interface{}) {
	fnVal := reflect.ValueOf(handlerFunc)
	fnType := fnVal.Type()

	// 校验签名是否符合：func(context.Context, Struct) error
	if fnType.Kind() != reflect.Func || fnType.NumIn() != 2 || fnType.NumOut() != 1 {
		panic("处理函数签名应为 func(context.Context, Struct) error")
	}

	// 获取主要入参
	payloadType := fnType.In(1)

	p.taskHandler[taskName] = func(ctx context.Context, payload any) error {
		val := reflect.New(payloadType).Interface()
		b, _ := json.Marshal(payload)
		if err := json.Unmarshal(b, val); err != nil {
			return fmt.Errorf("反序列化失败: %w", err)
		}

		// 调用函数  handlerFunc(ctx, payload)
		res := fnVal.Call([]reflect.Value{
			reflect.ValueOf(ctx), //context.Context 转换成 reflect.Value
			reflect.ValueOf(reflect.ValueOf(val).Elem().Interface()),
		})

		if errVal := res[0]; !errVal.IsNil() {
			return errVal.Interface().(error)
		}
		return nil
	}

	p.mux.HandleFunc(taskName, p.process)
}

// process 是统一处理入口
func (p *TaskProcessor) process(ctx context.Context, task *asynq.Task) error {
	handler, ok := p.taskHandler[task.Type()]
	if !ok {
		return fmt.Errorf("未注册任务处理器: %s", task.Type())
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("解析任务数据失败: %w", err)
	}

	return handler(ctx, payload)
}

// Start 启动任务处理服务
func (p *TaskProcessor) Start() error {
	if p.server == nil {
		return errors.New("任务服务未初始化")
	}
	return p.server.Start(p.mux)
}
