package publisher

import "context"

type Publisher interface {
	UploadVideo(ctx context.Context, filePath, title, description, keywords string, userId int) (string, error)
	Platform() string // 获取平台名称
	//RefrePlatformToken() error
}
