package mq

import (
	"errors"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// 定义rabbitmq的初始化，及数据推拉function
// RabbitMQ 实现了crawler.RabbitMQClient接口
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	redis     *cache.Redis
	queueName string
}

func NewRabbitConn() (*RabbitMQ, error) {
	config, err := config.LoadConfig("../../")
	if err != nil {
		return nil, errors.New("get config failed : " + err.Error())
	}
	rabbitMqConn, err := amqp.Dial(config.RabbitMQSource)
	if err != nil {
		return nil, errors.New("connect to rabbitmq failed : " + err.Error())
	}

	ch, err := rabbitMqConn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    rabbitMqConn,
		channel: ch,
		redis:   cache.RedisClient,
	}, nil
}
