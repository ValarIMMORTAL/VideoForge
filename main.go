package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"time"
)

type HotSearchResponse struct {
	Data struct {
		WordList []struct {
			Word       string `json:"word"`        // 热搜标题
			Position   int    `json:"position"`    // 排名
			HotValue   int    `json:"hot_value"`   // 热度值
			SentenceID string `json:"sentence_id"` // 句子ID（用于生成链接）
		} `json:"word_list"`
	} `json:"data"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.douyin.com"),
	)

	// 设置请求头（关键字段需更新为当前有效值）
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Host", "www.douyin.com")
		r.Headers.Set("Referer", "https://www.douyin.com/discover")
		r.Headers.Set("Cookie", "_xsrf=your_xsrf_token; _zap=your_zap_token") // 需替换为实际 Cookie
	})

	// 设置速率限制
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      2 * time.Second,
	})

	// 调用热搜接口
	apiURL := "https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web"
	c.OnResponse(func(r *colly.Response) {
		var resp HotSearchResponse
		if err := json.Unmarshal(r.Body, &resp); err != nil {
			fmt.Println("JSON 解析失败:", err)
			return
		}

		// 提取热搜数据
		for _, item := range resp.Data.WordList {
			fmt.Printf("排名: %d | 标题: %s | 热度: %d | 链接: https://www.douyin.com/hot/%s\n",
				item.Position, item.Word, item.HotValue, item.SentenceID)
		}
	})

	// 发送请求
	if err := c.Visit(apiURL); err != nil {
		fmt.Println("请求失败:", err)
	}
}
