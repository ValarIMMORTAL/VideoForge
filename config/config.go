package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	RabbitMQSource   string        `json:"RABBITMQSOURCE"`
	DouYingQueueName string        `json:"DOUYINGQUEUENAME"`
	MaxRetries       int           `json:"MAXRETRIES"`
	TimeOut          time.Duration `json:"TIMEOUT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	fmt.Println(path)
	viper.SetConfigName("app")
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
