package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/config"
	"github.com/pule1234/VideoForge/mq"
	"testing"
	"time"

	"github.com/pule1234/VideoForge/internal/models"
)

func TestPublish(t *testing.T) {
	// 创建测试数据
	testItem := models.TrendingItem{
		Title:     "测试标题",
		URL:       "http://test.com/video",
		ViewCount: "1000次观看",
		CreateAt:  time.Now().Truncate(time.Second), // 截断到秒，避免精度问题
	}
	loadConfig, err := config.LoadConfig("../../")
	if err != nil {
		fmt.Println(err)
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
	conn, err := mq.NewRabbitConn()
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.ConsumeItem(handler, loadConfig.DouYingQueueName)

}

func handler(item models.TrendingItem) error {
	res, _ := json.Marshal(item)
	fmt.Println(res)
	return nil
}
