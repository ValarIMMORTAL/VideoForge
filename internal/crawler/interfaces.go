package crawler

import (
	"github.com/gocolly/colly"
	"github.com/pule1234/VideoForge/internal/models"
	"time"
)

// Collector 定义爬虫收集器接口
type Collector interface {
	Visit(url string) error
	Wait()
}

// RabbitMQClient 定义RabbitMQ客户端接口
type RabbitMQClient interface {
	PublishItem(item models.TrendingItem) error
	CloseWithTimeout(timeout time.Duration) error
}

// 确保colly.Collector实现了Collector接口
var _ Collector = (*colly.Collector)(nil)
