package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/token"
)

func GetUserByToken(c *gin.Context, authorizationPayloadKey string) (int64, string, error) {
	payload, exists := c.Get(authorizationPayloadKey)
	if !exists {
		return 0, "", errors.New("authorization payload not found")
	}

	// 类型断言，将 interface{} 转换为具体类型
	authPayload, ok := payload.(*token.Payload)
	if !ok {
		return 0, "", errors.New("invalid authorization payload")
	}

	userID := authPayload.UserId
	userName := authPayload.Username

	return userID, userName, nil
}
