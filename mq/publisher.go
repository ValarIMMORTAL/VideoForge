package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

func (r *RabbitMQ) PublishItem(item interface{}, queueName string) error {
	//为当前消息生成唯一id
	msgId := uuid.New().String()
	message, _ := json.Marshal(item)

	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		queueName, true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "dlx.exchange",
			"x-dead-letter-routing-key": "dlx." + queueName,
		},
	)
	if err != nil {
		return fmt.Errorf("declare queue failed: %w", err)
	}

	// 开启确认模式
	if err = r.channel.Confirm(false); err != nil {
		return fmt.Errorf("enable confirm mode failed: %w", err)
	}

	confirmChan := r.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err = r.channel.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
			MessageId:   msgId,
		}); err != nil {
		return errors.New("push message failed : " + err.Error())
	}

	if err = r.waitForPendingMessages(confirmChan, msgId, message, queueName); err != nil {
		return err
	}

	log.Info().Msg("push message success queueName : " + queueName)
	return nil
}

// 确认消息送达
func (r *RabbitMQ) waitForPendingMessages(
	confirmChan chan amqp.Confirmation,
	msgId string,
	message []byte,
	queueName string,
) error {
	//设置超时机制
	timeout := time.After(10 * time.Second)
	for {
		select {
		case confirm := <-confirmChan:
			if confirm.Ack {
				log.Info().Msg("消息已确认 消息标识 ")
				r.redis.HSet(context.Background(), "msg:"+msgId, "status", "confirmed")
				return nil
			} else {
				log.Warn().Msg("消息未确认")
				return r.addRetryMessage(msgId, message, queueName)
			}
		case <-timeout:
			log.Warn().Msg("等待消息确认超时")
			return r.addRetryMessage(msgId, message, queueName)
		}
	}
}

func (r *RabbitMQ) addRetryMessage(msgId string, message []byte, queueName string) error {
	now := time.Now().Unix()
	nextRetry := now + 10 // 10秒后重试，可根据策略调整

	//开启管道， 将第一次发送失败的消息存储在hash 以及重试队列中
	pipe := r.redis.Conn.TxPipeline()
	pipe.HSet(context.Background(), "msg:"+msgId,
		"body", string(message),
		"status", "pending",
		"retry_count", 0,
		"next_retry", nextRetry,
	)
	pipe.ZAdd(context.Background(), "retry:"+queueName, redis.Z{Score: float64(nextRetry), Member: msgId})
	_, err := pipe.Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("写入重试队列失败")
	}
	return err
}

// 发送消息重试逻辑
func (r *RabbitMQ) StartRetryScheduler(interval time.Duration, maxRetry int) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.scanAndRetryQueues(maxRetry)
			}
		}
	}()
}

func (r *RabbitMQ) scanAndRetryQueues(maxRetry int) {
	ctx := context.Background()
	var cursor uint64
	for {
		// 查找所有 retry:* 的 ZSet key
		keys, newCursor, err := r.redis.Scan(ctx, cursor, "retry:*", 100)
		if err != nil {
			log.Error().Err(err).Msg("scan redis keys error")
			break
		}
		cursor = newCursor
		for _, key := range keys {
			queueName := strings.TrimPrefix(key, "retry:")
			r.retrySendMessages(ctx, queueName, maxRetry)
		}
		if cursor == 0 {
			break
		}
	}
}

func (r *RabbitMQ) retrySendMessages(ctx context.Context, queueName string, maxRetry int) {
	now := float64(time.Now().Unix())
	msgs, err := r.redis.ZRangeByScore(ctx, "retry:"+queueName, &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%.0f", now),
	})
	if err != nil {
		log.Error().Err(err).Msg("拉取重试任务失败")
		return
	}

	for _, msgId := range msgs {
		vals, err := r.redis.HGetAll(ctx, "msg:"+msgId)
		if err != nil || vals["status"] != "pending" { // 取数据失败 ｜｜ 数据不是正在等待处理Œ
			continue
		}

		retryCount, _ := strconv.Atoi(vals["retry_count"])
		if retryCount >= maxRetry {
			r.redis.HSet(ctx, "msg:"+msgId, "status", "failed")
		}

		err = r.PublishItem(json.RawMessage(vals["body"]), queueName)
		if err != nil {
			log.Warn().Err(err).Str("msgId", msgId).Msg("重新投递失败，稍后再试")
			r.redis.HSet(ctx, "msg:"+msgId, "retry_count", retryCount+1)
			continue
		}

		//处理成功  开启通道 删除 hash（msg: msgId）  删除zset "retry:"+queueName中Member = msgId
		pipe := r.redis.Conn.TxPipeline()
		pipe.HDel(ctx, "msg:"+msgId)
		pipe.ZRem(ctx, "retry:"+queueName, msgId)
		_, err = pipe.Exec(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("删除重试队列数据失败")
		}
		return
	}
}
