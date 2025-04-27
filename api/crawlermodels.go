package api

type DycrawlerRequest struct {
	Url string `json:"url" binding:"required"`
}
