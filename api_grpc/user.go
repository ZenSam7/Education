package api_grpc

import (
	"context"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/pb"
	"github.com/ZenSam7/Education/tools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (server *Server) mustEmbedUnimplementedEducationServer() {
	//TODO implement me
	panic("implement me")
}

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	passwordHash, err := tools.GetPasswordHash(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при хешировании: %s", err)
	}

	// Создаём пользователя
	arg := db.CreateUserParams{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: passwordHash,
	}
	user, err := server.queries.CreateUser(ctx, arg)
	if err != nil {
		// Если пользователь с таким именем уже есть, то выдаем ошибку
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			return nil, status.Errorf(codes.Internal, "пользователь с таким именем или email уже существует")
		}

		return nil, status.Errorf(codes.Internal, "не получилось создать пользователя: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convUser(user),
	}
	return response, nil
}
