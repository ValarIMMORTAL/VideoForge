package gapi

import (
	"context"
	"github.com/pule1234/VideoForge/global"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GenerateVideo(ctx context.Context, req *pb.GenerateVideoRequest) (resp *pb.GenerateVideoResponse, err error) {
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "%v", err)
	}

	userId := payload.UserId
	userName := payload.Username
	arg := processor.VideoParams{
		VideoSubject:        req.VideoSubject,
		VideoScript:         req.VideoScript,
		VideoTerms:          req.VideoTerms,
		VideoAspect:         req.VideoAspect,
		VideoConCatMode:     req.VideoConCatMode,
		VideoTransitionMode: req.VideoTransitionMode,
		VideoClipDuration:   int(req.VideoClipDuration),
		VideoCount:          int(req.VideoCount),
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
		FontSize:            int(req.FontSize),
		StrokeColor:         req.StrokeColor,
		StrokeWidth:         req.StrokeWidth,
		NThreads:            int(req.NThreads),
		ParagraphNumber:     int(req.ParagraphNumber),
	}

	//转换为context
	taskId, err := processor.GenerateVideo(global.GlobalCtx, arg, req.FileName, userName, userId, server.redis, server.store, server.qnManager, server.mq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	server.redis.SAdd(ctx, userName, taskId)
	resp.TaskId = taskId
	return
}
func (server *Server) GetVideos(ctx context.Context, req *pb.GetVideosRequest) (resp *pb.GetVideosResponse, err error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "请求参数无效: %v", err)
	}
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "%v", err)
	}

	userId := payload.UserId

	videos, err := server.store.GetVideosByUid(ctx, userId)
	if err != nil {
		return
	}
	temp := []*pb.Video{}
	for _, v := range videos {
		temp = append(temp, &pb.Video{
			Id:       v.ID,
			Title:    v.Title,
			Url:      v.Url,
			Duration: v.Duration,
		})
	}

	return &pb.GetVideosResponse{
		Videos: temp,
	}, nil
}

func validateGenerateVideoRequest(req *pb.GenerateVideoRequest) (violiations []*errdetails.BadRequest_FieldViolation) {

	return violiations
}

func validateGetVideosRequest(req *pb.GenerateVideoRequest) (violiations []*errdetails.BadRequest_FieldViolation) {

	return violiations
}
