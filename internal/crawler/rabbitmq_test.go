package crawler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/mq"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	// 创建测试数据
	item := processor.VideoMsg{
		Event:   "GenerateVideo failed",
		User_id: 1,
	}
	loadConfig, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Println("loadConfig err" + err.Error())
		return
	}

	//str, _ := json.Marshal(loadConfig)
	//fmt.Printf("loadconfig : %s", str)

	conn, err := mq.NewRabbitConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.PublishItem(item, loadConfig.DouYingQueueName)
	fmt.Println("发送成功！")
}

func TestConsumer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loadConfig, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Println(err)
		return
	}
	mq.InitRabbitMQ()

	conn, err := pgx.Connect(context.Background(), loadConfig.DBSource)
	if err != nil {
		t.Error("connect postgres err : " + err.Error())
	}
	q := db.New(conn)

	crawler, err := NewDyCrawler(loadConfig.DouYingQueueName, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	crawler.Rabbit.ConsumeItem(processor.CreateCopyWriting, loadConfig.DouYingQueueName, crawler.Postgres, ctx)

	time.Sleep(10 * time.Second)
}
