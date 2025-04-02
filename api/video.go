package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/global"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/token"
	"net/http"
)

// 视频处理相关接口
func (server *Server) generateVideo(c *gin.Context) {
	var req generateVideo

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//todo 进行用户认证，从token中获取到用户信息，将用户和视频生成的taskId进行绑定，存储在redis中，方便后续通知对应用户对应的视频已经生成
	payload, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	authPayload, ok := payload.(*token.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid authorization payload"})
		return
	}
	userName := authPayload.Username
	userId := authPayload.UserId
	arg := processor.VideoParams{
		VideoSubject:        req.VideoSubject,
		VideoScript:         req.VideoScript,
		VideoTerms:          req.VideoTerms,
		VideoAspect:         req.VideoAspect,
		VideoConCatMode:     req.VideoConCatMode,
		VideoTransitionMode: req.VideoTransitionMode,
		VideoClipDuration:   req.VideoClipDuration,
		VideoCount:          req.VideoCount,
		VideoSource:         req.VideoSource,
		VideoMaterals:       req.VideoMaterals,
		VideoLanguage:       req.VideoLanguage,
		VoiceName:           req.VoiceName,
		VoiceVolume:         req.VoiceVolume,
		VoiceRate:           req.VoiceRate,
		BgmType:             req.BgmType,
		BgmFile:             req.BgmFile,
		BgmVolume:           req.BgmVolume,
		SubtitleEnabled:     req.SubtitleEnabled,
		SubtitlePosition:    req.SubtitlePosition,
		CustomPosition:      req.CustomPosition,
		FontName:            req.FontName,
		TextForeColor:       req.TextForeColor,
		TextBackgroundColor: req.TextBackgroundColor,
		FontSize:            req.FontSize,
		StrokeColor:         req.StrokeColor,
		StrokeWidth:         req.StrokeWidth,
		NThreads:            req.NThreads,
		ParagraphNumber:     req.ParagraphNumber,
	}

	//转换为context
	result, err := processor.GenerateVideo(global.GlobalCtx, arg, userName, userId, server.redis, server.store, server.qnManager)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//将返回的taskID和user_id存储到数据库中
	server.redis.SAdd(c, userName, result)

	c.JSON(http.StatusOK, result)
}
