package gapi

import (
	"context"
	"errors"
	"github.com/lib/pq"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/pb"
	"github.com/pule1234/VideoForge/util"
	"github.com/pule1234/VideoForge/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (resp *pb.CreateUserResponse, err error) {
	violiations := validateCreateuser(req)
	if violiations != nil {
		return nil, invalidArgumentError(violiations)
	}
	password, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.UserName,
		HashedPassword: password,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %s", pqErr.Detail)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	resp = &pb.CreateUserResponse{
		UserId:   user.ID,
		UserName: user.Username,
		Email:    user.Email,
	}

	return resp, nil
}

func (server *Server) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {
	violiations := validateUserLogin(req)
	if violiations != nil {
		return nil, invalidArgumentError(violiations)
	}
	user, err := server.store.GetUserByName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "userName not exist: %s", req.UserName)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	//check password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	//create accesstoken
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to createToken: %v", err)
	}

	//create refreToken
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to createToken: %v", err)
	}
	mtdt := server.Metadata(ctx)
	arg := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}
	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to createSession: %v", err)
	}

	res := &pb.UserLoginResponse{
		User:                  converter(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return res, nil
}

func validateCreateuser(req *pb.CreateUserRequest) (violiations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateUsername(req.UserName)
	if err != nil {
		violiations = append(violiations, fieldViolation("username", err))
	}
	err = val.ValidatePassword(req.Password)
	if err != nil {
		violiations = append(violiations, fieldViolation("password", err))
	}
	err = val.ValidateEmail(req.Email)
	if err != nil {
		violiations = append(violiations, fieldViolation("email", err))
	}

	return violiations
}

func validateUserLogin(req *pb.UserLoginRequest) (violiations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateUsername(req.UserName)
	if err != nil {
		violiations = append(violiations, fieldViolation("username", err))
	}
	err = val.ValidatePassword(req.Password)
	if err != nil {
		violiations = append(violiations, fieldViolation("password", err))
	}

	return violiations
}
