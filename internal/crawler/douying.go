package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/models"
	"github.com/pule1234/VideoForge/mq"
	"github.com/rs/zerolog/log"
)

type DyCrawler struct {
	Collector Collector
	QueueName string
	Rabbit    *mq.RabbitMQ //全局mq连接
	Postgres  db.Store     // 数据库连接， 所有的Crawler都会携带一个连接， processor便不需要连接
}

func NewDyCrawler(queueName string, queries db.Store) (*DyCrawler, error) {
	baseCrawler, err := NewCrawler() //已经定义好错误处理
	if err != nil {
		return nil, fmt.Errorf("创建爬虫失败: %v", err)
	}

	dc := &DyCrawler{
		Collector: baseCrawler,
		QueueName: queueName,
		Rabbit:    mq.GlobalRabbitMQ,
		Postgres:  queries,
	}

	// 设置爬虫回调
	baseCollector := baseCrawler.collector

	baseCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Host", "www.douyin.com")
		r.Headers.Set("Referer", "https://www.douyin.com/discover")
		r.Headers.Set("Cookie", "_xsrf=your_xsrf_token; _zap=your_zap_token") //当前接口为免登录接口
	})

	baseCollector.OnResponse(func(r *colly.Response) {
		var resp HotSearchResponse
		if err := json.Unmarshal(r.Body, &resp); err != nil {
			fmt.Println("JSON 解析失败:", err)
			return
		}

		var mqItem []models.TrendingItem
		// 提取热搜数据
		for _, data := range resp.Data.WordList {
			fmt.Printf("排名: %d | 标题: %s | 热度: %d | 链接: https://www.douyin.com/hot/%s\n",
				data.Position, data.Word, data.HotValue, data.SentenceID)
			// 修改 CreateCopyWriting  一次性处理多个item （[]models.TrendingItem）   接受crawler的数据库连接
			item := models.TrendingItem{
				Source:   "DouYing",
				Title:    data.Word,
				Position: data.Position,
			}
			mqItem = append(mqItem, item)
		}
		loadConfig, _ := config.LoadConfig("../../")
		if err = dc.Rabbit.PublishItem(mqItem, loadConfig.DouYingQueueName); err != nil {
			log.Error().Err(err).Msg("爬虫推送数据失败")
			return
		}
	})

	return dc, nil
}

// url := "https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web"
func (d *DyCrawler) Start(url string) error {
	err := d.Collector.Visit(url)
	if err != nil {
		return errors.New("DouYing_Crawler error: " + err.Error() + "\n DouYing_URL: " + url)
	}

	return nil
}

// 优化关闭方法
//func (d *DyCrawler) Stop() {
//	if d.Rabbit != nil {
//		if err := d.Rabbit.CloseWithTimeout(10 * time.Second); err != nil { // 延长关闭超时时间
//			log.Printf("RabbitMQ关闭错误: %v", err)
//		}
//	}
//}
