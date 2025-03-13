package models

import "time"

type TrendingItem struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	ViewCount string    `json:"view_count"`
	CreateAt  time.Time `json:"create_at"`
	Source    string    `json:"source"` //数据来源，方便将生成的文案对应到视频平台
}
