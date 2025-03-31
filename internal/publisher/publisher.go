package publisher

import (
	"context"
	db "github.com/pule1234/VideoForge/db/sqlc"
)

type Publisher interface {
	UploadVideo(ctx context.Context, filePath, title, description, keywords string, userId int64, store db.Store) (string, error)
	Platform() string // 获取平台名称
	//RefrePlatformToken() error
}
