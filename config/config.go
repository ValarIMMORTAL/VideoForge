package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	RabbitMQSource       string        `mapstructure:"RABBITMQSOURCE"`
	DouYingQueueName     string        `mapstructure:"DOUYINGQUEUENAME"`
	MaxRetries           int           `mapstructure:"MAXRETRIES"`
	TimeOut              time.Duration `mapstructure:"TIMEOUT"`
	AiUrl                string        `mapstructure:"AIURL"`
	ApiKey               string        `mapstructure:"APIKEY"`
	AiModel              string        `mapstructure:"AIMODEL"`
	Role                 string        `mapstructure:ROLE`
	CopyWritingContent   string        `mapstructure:COPYWRITINGCONTENT`
	GenerateVideoBaseUrl string        `mapstructure:GENERATEVIDEOBASEURL`
	VideoEndpoint        string        `mapstructure:VIDEOENDPOINT`
	TaskEndpoint         string        `mapstructure:TASKENDPOINT`
	RedisSource          string        `mapstructure:"REDISSOURCE"`
	RedisPassword        string        `mapstructure:"REDISPASSWORD"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	VideoPath            string        `mapstructure:"VIDEOPATH"`
	CdnDomain            string        `mapstructure:"CDNDOMAIN"`
	TempDir              string        `mapstructure:"TEMPDIR"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	configName := os.Getenv("APP_ENV")
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	//自动检查环境
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config *Config
	err = viper.Unmarshal(&config)
	return config, err
}
