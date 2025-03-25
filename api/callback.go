package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/global"
	"net/http"
)

func (server *Server) callback(ctx *gin.Context) {
	code := ctx.Query("code") // 从 URL 参数获取 code
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'code' parameter"})
		return
	}
	fmt.Println("ping 中的 code为 : " + code)
	// 非阻塞发送 code（避免通道未就绪时阻塞）
	select {
	case global.OauthCodeChan <- code:
		ctx.String(http.StatusOK, "OAuth code received. You can close this window.")
	default:
		ctx.String(http.StatusOK, "OAuth code already processed.")
	}
}
