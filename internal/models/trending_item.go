package models

import "time"

type TrendingItem struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	ViewCount string    `json:"view_count"`
	CreateAt  time.Time `json:"create_at"`
}