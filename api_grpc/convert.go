package api_grpc

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convUser(user db.User) *pb.User {
	return &pb.User{
		IdUser:      user.IDUser,
		Name:        user.Name,
		Description: user.Description.String,
		Email:       user.Email,
		Karma:       user.Karma,
		CreatedAt:   timestamppb.New(user.CreatedAt.Time),
	}
}
