package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pule1234/VideoForge/global"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/util"
	"net/http"
)

// 视频处理相关接口
func (server *Server) generateVideo(c *gin.Context) {
	var req generateVideo

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, userName, err := util.GetUserByToken(c, authorizationPayloadKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	result, err := processor.GenerateVideo(global.GlobalCtx, arg, req.FileName, userName, userId, server.redis, server.store, server.qnManager, server.mq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//将返回的taskID和user_id存储到数据库中
	server.redis.SAdd(c, userName, result)

	c.JSON(http.StatusOK, gin.H{
		"task_id": result,
	})
}

// 获取用户视频列表
func (server *Server) getVideos(c *gin.Context) {
	var req getVideosByUidReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID, _, err := util.GetUserByToken(c, authorizationPayloadKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error get user_id failed": err.Error()})
		return
	}
	videos, err := server.store.GetVideosByUid(c, userID)
	if err != nil {
		return
	}
	temp := []Videos{}
	for _, video := range videos {
		temp = append(temp, Videos{
			id:    video.ID,
			title: video.Title,
			url:   video.Url,
		})
	}
	resp := getVideosByUidResp{
		Videos: temp,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}
