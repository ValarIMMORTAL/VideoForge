package mq

import "sync"

var (
	GlobalRabbitMQ *RabbitMQ // 全局使用rabbitmq连接
	once           sync.Once // 确保只初始化一次
)

// main中调用此函数，初始化全局rabbitmq变量
func InitRabbitMQ() error {
	var err error
	once.Do(func() {
		GlobalRabbitMQ, err = NewRabbitConn()
	})

	return err
}

// GetRabbitMQ 获取全局 RabbitMQ 实例
func GetRabbitMQ() *RabbitMQ {
	return GlobalRabbitMQ
}
