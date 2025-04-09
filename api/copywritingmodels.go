package api

import (
	"time"
)

type getCopyWritingRequest struct {
	date time.Time `json:"date" binding:"required"`
	page int32     `json:"page" binding:"required"`
	num  int32     `json:"num" binding:"required"`
}

type getCopyWritingResponse struct {
	items []Copywriting `json:"items"`
}

type Copywriting struct {
	ID        int64     `json:"id"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	DeleteAt  time.Time `json:"delete_at"`
}
