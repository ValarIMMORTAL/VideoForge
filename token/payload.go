package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// 自定义token错误返回
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// 存储token中的有效数据
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserId    int32     `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`  // 创建时间
	ExpiredAt time.Time `json:"expired_at"` //过期时间
}

func NewPayload(userId int32, username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        tokenID,
		UserId:    userId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// 判断是否过期
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
