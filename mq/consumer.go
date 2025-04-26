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
	err := r.channel.ExchangeDeclare(
		"dlx.exchange",
		"direct",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("declare dead letter exchange failed")
		return
	}

	q, err := r.channel.QueueDeclare(
		queueName, true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "dlx.exchange",
			"x-dead-letter-routing-key": "dlx." + queueName,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("declare queue failed")
		return
	}

	msgs, err := r.channel.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("consume failed")
		return
	}

	const maxRetry = 3 // 设置最大重试次数

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
					_ = d.Reject(false)
					continue
				}

				if err := handler(items, dbStore); err != nil {
					log.Error().Err(err).Msg("处理失败")

					retryCount := 0
					if rc, ok := d.Headers["x-retry-count"].(int32); ok {
						retryCount = int(rc)
					}

					if retryCount >= maxRetry {
						log.Error().Int("retry", retryCount).Msg("超过最大重试次数，进入死信队列")
						_ = d.Reject(false)
					} else {
						log.Warn().Int("retry", retryCount).Msg("处理失败，重新投递")
						_ = d.Reject(false)
						// 重新发布消息（带上retry+1）
						r.channel.Publish(
							"", queueName, false, false,
							amqp.Publishing{
								ContentType: "application/json",
								Body:        d.Body,
								Headers: amqp.Table{
									"x-retry-count": retryCount + 1,
								},
							},
						)
					}
					continue
				}

				//处理成功，ack
				_ = d.Ack(false)
			}
		}
	}()
}
