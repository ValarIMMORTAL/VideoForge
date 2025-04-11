package gapi

import (
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/pb"
)

func converter(user db.User) *pb.User {
	return &pb.User{
		Userid:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
