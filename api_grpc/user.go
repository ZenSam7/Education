package api_grpc

import (
	"context"
	"database/sql"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/ZenSam7/Education/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func validateCreateUserRequest(req *pb.CreateUserRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := validator.ValidateString(req.GetName(), 1, 99); err != nil {
		errorsFields = append(errorsFields, fieldViolation("name", err))
	}

	if err := validator.ValidateString(req.GetPassword(), 1, 999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("password", err))
	}

	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		errorsFields = append(errorsFields, fieldViolation("email", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

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
			return nil, status.Errorf(codes.AlreadyExists, "пользователь с таким именем или email уже существует")
		}
		return nil, status.Errorf(codes.Internal, "не получилось создать пользователя: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convUser(user),
	}
	return response, nil
}

func validateGetUserRequest(req *pb.GetUserRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := validator.ValidateNaturalNum(int(req.GetIdUser())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_user", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := validateGetUserRequest(req); err != nil {
		return nil, err
	}

	user, err := server.queries.GetUser(ctx, req.GetIdUser())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить пользователя: %s", err)
	}

	response := &pb.GetUserResponse{
		User: convUser(user),
	}
	return response, nil
}

func validateGetManySortedUsersRequest(req *pb.GetManySortedUsersRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := validator.ValidateNaturalNum(int(req.GetPageNum())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_num", err))
	}

	if err := validator.ValidateNaturalNum(int(req.GetPageSize())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_size", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) GetManySortedUsers(ctx context.Context, req *pb.GetManySortedUsersRequest) (*pb.GetManySortedUsersResponse, error) {
	if err := validateGetManySortedUsersRequest(req); err != nil {
		return nil, err
	}

	arg := db.GetManySortedUsersParams{
		IDUser:      req.GetIdUser(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Karma:       req.GetKarma(),
		Limit:       req.GetPageSize(),
		Offset:      (req.GetPageNum() - 1) * req.GetPageSize(),
	}
	users, err := server.queries.GetManySortedUsers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "пользователи не найдены")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить пользователей: %s", err)
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, convUser(u))
	}

	response := &pb.GetManySortedUsersResponse{
		Users: pbUsers,
	}
	return response, nil
}

func validateEditUserRequest(req *pb.EditUserRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if req.Name != nil {
		if err := validator.ValidateString(req.GetName(), 1, 0); err != nil {
			errorsFields = append(errorsFields, fieldViolation("page_num", err))
		}
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) EditUser(ctx context.Context, req *pb.EditUserRequest) (*pb.EditUserResponse, error) {
	if err := validateEditUserRequest(req); err != nil {
		return nil, err
	}

	accessPayload, err := server.authUser(ctx)
	if err != nil {
		return nil, err
	}

	arg := db.EditUserParams{
		IDUser: accessPayload.IDUser,
		Name:   req.GetName(),
		// Разделяем пустое значение и значение которое вообще не указывали
		// (Т.е. имеем возможность указать '' или 0 как валидный параметр)
		Description: pgtype.Text{String: req.GetDescription(), Valid: req.Description != nil},
		Karma:       pgtype.Int4{Int32: req.GetKarma(), Valid: req.Karma != nil},
	}

	editedUser, err := server.queries.EditUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось изменить пользователя: %s", err)
	}

	response := &pb.EditUserResponse{
		User: convUser(editedUser),
	}
	return response, nil
}

func (server *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_ = req.String() // Просто чтобы не было предупреждений

	payload, err := server.authUser(ctx)
	if err != nil {
		return nil, err
	}

	deletedUser, err := server.queries.DeleteUser(ctx, payload.IDUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось удалить пользователя: %s", err)
	}

	response := &pb.DeleteUserResponse{
		User: convUser(deletedUser),
	}
	return response, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := validator.ValidateString(req.GetPassword(), 1, 999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("password", err))
	}

	if err := validator.ValidateString(req.GetName(), 1, 99); err != nil {
		errorsFields = append(errorsFields, fieldViolation("name", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if err := validateLoginUserRequest(req); err != nil {
		return nil, err
	}

	user, err := server.queries.GetUserFromName(ctx, req.GetName())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить пользователя: %s", err)
	}

	if !tools.CheckPassword(req.GetPassword(), user.PasswordHash) {
		return nil, status.Errorf(codes.Unauthenticated, "неправильный пароль")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.IDUser, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать access token: %s", err)
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.IDUser, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать refresh token: %s", err)
	}

	info := server.extractMetadata(ctx)
	_, err = server.queries.CreateSession(ctx, db.CreateSessionParams{
		IDSession:    pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true},
		IDUser:       user.IDUser,
		RefreshToken: refreshToken,
		ExpiredAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiredAt, Valid: true},
		ClientIp:     info.ClientIP,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать сессию: %s", err)
	}

	response := &pb.LoginUserResponse{
		User:                  convUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return response, nil
}

func validateRenewAccessTokenRequest(req *pb.RenewAccessTokenRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := validator.ValidateString(req.GetRefreshToken(), 1, 9999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("password", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	if err := validateRenewAccessTokenRequest(req); err != nil {
		return nil, err
	}

	info := server.extractMetadata(ctx)

	refreshPayload, errVerifyToken := server.tokenMaker.VerifyToken(req.GetRefreshToken())

	oldSession, err := server.queries.DeleteSession(ctx, pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось удалить сессию: %s", err)
	} else if oldSession.Blocked {
		return nil, status.Errorf(codes.Unauthenticated, "сессия заблокирована")
	} else if oldSession.IDUser != refreshPayload.IDUser {
		return nil, status.Errorf(codes.Unauthenticated, "некорректная сессия пользователя")
	} else if oldSession.RefreshToken != req.GetRefreshToken() {
		return nil, status.Errorf(codes.Unauthenticated, "некорректный refresh token")
	}
	if errVerifyToken == token.ErrorExpiredToken || oldSession.ClientIp != info.ClientIP {
		return nil, status.Errorf(codes.Unauthenticated, "необходимо залогиниться")
	}

	newRefreshToken, newRefreshPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать новый refresh token: %s", err)
	}
	newAccessToken, newAccessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать новый access token: %s", err)
	}

	_, err = server.queries.CreateSession(ctx, db.CreateSessionParams{
		IDSession:    pgtype.UUID{Bytes: newRefreshPayload.IDSession, Valid: true},
		IDUser:       newRefreshPayload.IDUser,
		RefreshToken: newRefreshToken,
		ExpiredAt:    pgtype.Timestamptz{Time: newRefreshPayload.ExpiredAt, Valid: true},
		ClientIp:     info.ClientIP,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать новую сессию: %s", err)
	}

	response := &pb.RenewAccessTokenResponse{
		AccessToken:           newAccessToken,
		RefreshToken:          newRefreshToken,
		AccessTokenExpiredAt:  timestamppb.New(newAccessPayload.ExpiredAt),
		RefreshTokenExpiredAt: timestamppb.New(newRefreshPayload.ExpiredAt),
	}

	return response, nil
}
