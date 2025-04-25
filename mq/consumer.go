package mq

import (
	"context"
	"encoding/json"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// simple 模式下消费者
func (r *RabbitMQ) ConsumeItem(handler func(item []models.TrendingItem, dbStore *db.Queries) error, queueName string, dbStore *db.Queries, ctx context.Context) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		queueName, false, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "dlx.exchange",
			"x-dead-letter-routing-key": "dlx." + queueName,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("consumer queuedeclare failed")
		return
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("comsumer consume message err")
		return
	}

	//启用协程处理消息
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Received shutdown signal. Exiting...")
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				log.Printf("Received a message: %s", d.Body)

				var items []models.TrendingItem
				if err := json.Unmarshal(d.Body, &items); err != nil {
					log.Error().Err(err).Msg("unmarshal failed, reject message")
					_ = d.Reject(false) // 拒绝并不重回队列，进入 DLX
					continue
				}

				if err = handler(items, dbStore); err != nil {
					log.Error().Err(err).Msg("处理失败，消息进入死信队列")
					_ = d.Reject(false) // 不重回队列，进入死信队列
					continue
				}
				_ = d.Ack(false)
			}
		}
	}()

	return
}
