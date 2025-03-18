package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	"log"
	"time"

	"github.com/pule1234/VideoForge/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// 定义rabbitmq的初始化，及数据推拉function
// RabbitMQ 实现了crawler.RabbitMQClient接口
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitConn() (*RabbitMQ, error) {
	config, err := config.LoadConfig("../../")
	if err != nil {
		return nil, errors.New("get config failed : " + err.Error())
	}
	fmt.Println(config.RabbitMQSource)
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
	}, nil
}

func (r *RabbitMQ) PublishItem(item []models.TrendingItem, queueName string) error {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		queueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	message, _ := json.Marshal(item)
	//调用channel 发送消息到队列中
	err = r.channel.Publish(
		"",
		queueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return errors.New("push message failed : " + err.Error())
	}

	//if err = r.waitForPendingMessages(); err != nil {
	//	log.Println("等待消息发送超时:", err)
	//	return errors.New("等待消息发送超时" + err.Error())
	//}
	return nil
}

// simple 模式下消费者
// todo 将handler 替换成生成文案的functuon ， 并且将爬取到的关键字（keyword）和关键字来源（item.source）作为参数传入到function中
func (r *RabbitMQ) ConsumeItem(handler func(item []models.TrendingItem, dbStore *db.Queries) error, queueName string, dbStore *db.Queries) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		queueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		log.Println("consumer queuedeclare failed:", err.Error())
		return
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Println("comsumer consume message err : ", err.Error())
	}

	forever := make(chan bool)
	errChan := make(chan error)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)

			var items []models.TrendingItem
			err = json.Unmarshal(d.Body, &items)
			if err != nil {
				return
			}

			err = handler(items, dbStore) //关键字(title) 和信息来源(source)都在item中
			if err != nil {
				errChan <- err
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	err = <-errChan
	log.Printf("handler err is :" + err.Error())
	<-forever

}

// 新增带重试的连接方法
func NewRabbitConnWithRetry(maxRetries int) (*RabbitMQ, error) {
	var conn *RabbitMQ
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = NewRabbitConn() // 使用默认队列名，也可以从配置中获取
		if err == nil {
			return conn, nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return nil, fmt.Errorf("经过%d次重试后连接失败: %v", maxRetries, err)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// CloseWithTimeout 带超时的关闭
func (r *RabbitMQ) CloseWithTimeout(timeout time.Duration) error {
	done := make(chan bool)
	go func() {
		// 等待所有消息发送完成
		if err := r.waitForPendingMessages(); err != nil {
			log.Println("等待消息发送超时:", err)
		}
		// 关闭通道和连接
		if r.channel != nil {
			r.channel.Close()
		}
		if r.conn != nil {
			r.conn.Close()
		}
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.New("close timeout")
	}
}

// 确认消息送达
func (r *RabbitMQ) waitForPendingMessages() error {
	//启用 Publisher Confirms 模式。 告知消息是否发送成功
	confirmChan := r.channel.NotifyPublish(make(chan amqp.Confirmation, 10))

	//设置超时机制
	timeout := time.After(10 * time.Second)
	for {
		select {
		case confirm := <-confirmChan:
			if confirm.Ack {
				log.Println("消息已确认 消息表示 : ", confirm.DeliveryTag)
				return nil
			} else {
				log.Println("消息未确认")
				return errors.New("pending message failed")
			}
		case <-timeout:
			return errors.New("等待消息确认超时")
		}
	}
}
