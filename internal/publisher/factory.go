package publisher

import "fmt"

type PlatformConfig struct {
	Platform string                 //平台名称
	Config   map[string]interface{} //平台相关配置
}

type PublisherFactory struct {
	platformConfigs map[string]PlatformConfig
}

func NewPublisherFactory() *PublisherFactory {
	return &PublisherFactory{
		platformConfigs: loadPlatformConfigs(), //查询所有平台的配置
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
func loadPlatformConfigs() map[string]PlatformConfig {
	return map[string]PlatformConfig{
		"youtube": PlatformConfig{},
	}
}
