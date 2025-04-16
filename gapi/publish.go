package gapi

import (
	"context"
	"github.com/pule1234/VideoForge/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

// todo 从gin中接收文件
func (server *Server) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "请求参数无效: %v", err)
	}

	filePath, err := server.qnManager.DownloadFile(server.config.TempDir, req.FileName, req.UserName, req.Bucket, req.Subscribe, req.Domain)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "downloadFile failed : %v", err)
	}
	defer os.Remove(filePath)
	// 4. 调用 Publisher 上传
	//var userID int32 = 1
	publisher, err := server.publisherFactory.CreatePublisher(req.PlatformName)
	videoID, err := publisher.UploadVideo(ctx, filePath, req.Title, req.Description, "keyword", req.UserId, server.store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "upload videos failed : %v", err)
	}
	resp := &pb.UploadVideoResponse{
		VideoId: videoID,
	}
	return resp, nil
}
