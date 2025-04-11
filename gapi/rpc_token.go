package gapi

import (
	"context"
	"github.com/pule1234/VideoForge/pb"
	"github.com/pule1234/VideoForge/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// input RefreshToken then return a new AccessToken
func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (resp *pb.RenewAccessTokenResponse, err error) {
	violidations := validateReNewAccessToken(req)
	if violidations != nil {
		return nil, invalidArgumentError(violidations)
	}
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	if session.IsBlocked {
		return nil, status.Error(codes.Unauthenticated, "session is blocked")
	}

	if session.UserID != refreshPayload.UserId {
		return nil, status.Error(codes.PermissionDenied, "session does not have permission to renew token")
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, status.Error(codes.Unauthenticated, "refresh token invalid")
	}

	//验证当前refreshtoken是否过期
	if time.Now().After(refreshPayload.ExpiredAt) {
		return nil, status.Error(codes.PermissionDenied, "session expired")
	}

	//当前token有效生成accessToken
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserId, refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "can not create access token")
	}
	resp = &pb.RenewAccessTokenResponse{
		AccessToken:         accessToken,
		AccessTokenExpireAt: timestamppb.New(accessPayload.ExpiredAt),
	}

	return resp, nil
}

func validateReNewAccessToken(req *pb.RenewAccessTokenRequest) (violidations []*errdetails.BadRequest_FieldViolation) {
	// 校验refreshToken是否有效
	if err := val.ValidateRefreshToken(req.GetRefreshToken()); err != nil {
		violidations = append(violidations, fieldViolation("refreshToken", err))
	}
	return violidations
}
