package crawler

import (
	"github.com/pule1234/VideoForge/internal/models"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockRabbitMQ 用于在测试中模拟RabbitMQ连接
type MockRabbitMQ struct {
	mock.Mock
}

// PublishItem 模拟发布消息到队列
func (m *MockRabbitMQ) PublishItem(item models.TrendingItem, queueName string) error {
	args := m.Called(item, queueName)
	return args.Error(0)
}

// CloseWithTimeout 模拟带超时的关闭连接
func (m *MockRabbitMQ) CloseWithTimeout(timeout time.Duration) error {
	args := m.Called(timeout)
	return args.Error(0)
}