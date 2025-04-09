package gapi

import (
	"context"
	"github.com/pule1234/VideoForge/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sever *Server) RenewAccessToken(context.Context, *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method RenewAccessToken not implemented")
}
