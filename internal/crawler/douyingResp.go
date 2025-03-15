package crawler

// dy 热榜接口返回值
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
