package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pule1234/VideoForge/mq"

	"github.com/gocolly/colly"
	"github.com/pule1234/VideoForge/internal/models"
)

type DyCrawler struct {
	collector Collector
	queueName string
	rabbit    *mq.RabbitMQ
}

func newDyCrawler(queueName string) (*DyCrawler, error) {
	baseCrawler, err := NewCrawler() //已经定义好错误处理
	if err != nil {
		return nil, fmt.Errorf("创建爬虫失败: %v", err)
	}

	dc := &DyCrawler{
		collector: baseCrawler,
		queueName: queueName,
		rabbit:    nil,
	}

	// 设置爬虫回调
	baseCollector := baseCrawler.collector

	baseCollector.OnHTML("div.popular-item", func(e *colly.HTMLElement) {
		item := models.TrendingItem{
			Title:     e.ChildText("div.title"),
			URL:       e.ChildAttr("a", "href"),
			ViewCount: e.ChildText("span.view-count"),
			//Source:
			CreateAt: time.Now(),
		}
		res, _ := json.Marshal(&item)
		fmt.Println(res)
		// 发送到消息队列
		if dc.rabbit != nil {
			if err := dc.rabbit.PublishItem(item, queueName); err != nil {
				log.Printf("Failed to publish item: %v", err)
			}
		}
	})

	return dc, nil
}

func (d *DyCrawler) Start(url string) error {
	ticker := time.NewTicker(24 * time.Hour)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 增加连接超时控制
	defer cancel()

	go func() {
		for {
			var rabbitMQ *mq.RabbitMQ
			var err error

			// 使用带超时的连接方法
			select {
			case <-ctx.Done():
				log.Printf("RabbitMQ连接超时")
				return
			default:
				rabbitMQ, err = mq.NewRabbitConnWithRetry(5) // 增加重试次数
			}

			if err != nil {
				log.Printf("等待下一个周期...")
				<-ticker.C
				continue
			}

			d.rabbit = rabbitMQ

			// 执行爬取任务
			if err := d.collector.Visit(url); err != nil {
				log.Printf("爬取数据失败: %v", err)
			}

			// 关闭当前RabbitMQ连接
			d.Stop()
			// 等待下一次定时触发
			<-ticker.C
		}
	}()
	return nil
}

// 优化关闭方法
func (d *DyCrawler) Stop() {
	if d.rabbit != nil {
		if err := d.rabbit.CloseWithTimeout(10 * time.Second); err != nil { // 延长关闭超时时间
			log.Printf("RabbitMQ关闭错误: %v", err)
		}
	}
}
