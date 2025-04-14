package gapi

import (
	"context"
	"github.com/pule1234/VideoForge/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// todo 从gin中接收文件
func (server *Server) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "请求参数无效: %v", err)
	}

	return nil, status.Errorf(codes.Unimplemented, "method UploadVideo not implemented")
}
