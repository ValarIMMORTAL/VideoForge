package models

// 爬虫数据结构体
type TrendingItem struct {
	URL      string `json:"url"`
	Source   string `json:"source"`   //数据来源，方便将生成的文案对应到视频平台
	Title    string `json:"title"`    //爬取到的标题
	Position int    `json:"position"` //排名
}
