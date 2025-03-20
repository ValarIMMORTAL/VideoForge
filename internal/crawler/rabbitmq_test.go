package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/mq"
	"testing"
)

func TestPublish(t *testing.T) {
	// 创建测试数据
	testItem := []models.TrendingItem{}
	testItem = append(testItem, models.TrendingItem{
		Title: "测试标题",
		URL:   "http://test.com/video",
	})
	loadConfig, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Println("loadConfig err" + err.Error())
		return
	}

	str, _ := json.Marshal(loadConfig)
	fmt.Printf("loadconfig : %s", str)

	conn, err := mq.NewRabbitConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.PublishItem(testItem, loadConfig.DouYingQueueName)
	fmt.Println("发送成功！")
}

func TestConsumer(t *testing.T) {
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

	crawler, err := newDyCrawler(loadConfig.DouYingQueueName, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	crawler.rabbit.ConsumeItem(processor.CreateCopyWriting, loadConfig.DouYingQueueName, crawler.postgres)

}
