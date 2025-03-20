package config

import (
	"fmt"
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
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	configName := os.Getenv("APP_ENV")
	fmt.Println("app_env" + configName)
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
