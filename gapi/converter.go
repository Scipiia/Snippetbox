package gapi

import (
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Name:              user.Name,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		Created:           timestamppb.New(user.Created),
	}
}
