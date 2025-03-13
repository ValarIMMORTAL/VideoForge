package crawler

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/pule1234/VideoForge/config"
)

type Crawler struct {
	collector *colly.Collector
}

// Visit 实现Collector接口的Visit方法
func (c *Crawler) Visit(url string) error {
	return c.collector.Visit(url)
}

// Wait 实现Collector接口的Wait方法
func (c *Crawler) Wait() {
	c.collector.Wait()
}

func NewCrawler() (*Crawler, error) {
	config, _ := config.LoadConfig(".")
	c := colly.NewCollector(
		// 设置http请求的User_Agent 模拟浏览器访问
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"),
		// 爬取的最大深度
		colly.MaxDepth(2),
	)

	// 设置限制
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		retries := 0
		for retries < config.MaxRetries {
			if err := r.Request.Retry(); err == nil {
				return
			}
			retries++
			time.Sleep(time.Second * time.Duration(retries))
		}
	})

	return &Crawler{
		collector: c,
	}, nil
}
