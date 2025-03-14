package crawler

import (
	"github.com/gocolly/colly"
)

// Collector 定义爬虫收集器接口
type Collector interface {
	Visit(url string) error
	Wait()
}

// 确保colly.Collector实现了Collector接口
var _ Collector = (*colly.Collector)(nil)
