package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	db "github.com/pule1234/VideoForge/db/sqlc"
)

type PlatformConfig struct {
	Platform string                 //平台名称
	Config   map[string]interface{} //平台相关配置
}

type PublisherFactory struct {
	platformConfigs map[string]PlatformConfig
}

func NewPublisherFactory(store db.Store) *PublisherFactory {
	return &PublisherFactory{
		platformConfigs: loadPlatformConfigs(store), //查询所有平台的配置
	}
}

func (f *PublisherFactory) CreatePublisher(platformName string) (Publisher, error) {
	config, ok := f.platformConfigs[platformName]
	if !ok {
		return nil, fmt.Errorf("平台 %s 不存在", platformName)
	}
	switch platformName {
	case "youtube":
		return NewYouTubePublisher(config)
	default:
		return nil, fmt.Errorf("不支持的平台: %s", platformName)
	}
}

// 查询所有平台的配置
func loadPlatformConfigs(store db.Store) map[string]PlatformConfig {
	platforms, err := store.GetPlatforms(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	platformConfigs := map[string]PlatformConfig{}
	for _, platform := range platforms {
		var temp map[string]interface{}
		err = json.Unmarshal([]byte(platform.Detail), &temp)
		if err != nil {
			fmt.Println("JSON 解析失败:", err)
			return nil
		}
		platformConfigs[platform.Platform] = PlatformConfig{
			Platform: platform.Platform,
			Config:   temp,
		}
	}

	return platformConfigs
}
