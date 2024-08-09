package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	worker2 "github.com/ZenSam7/Education/redis/worker"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

func validateCreateUserRequest(req *pb.CreateUserRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateString(req.GetName(), 1, 99); err != nil {
		errorsFields = append(errorsFields, fieldViolation("name", err))
	}

	if err := tools.ValidateString(req.GetPassword(), 1, 999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("password", err))
	}

	if err := tools.ValidateEmail(req.GetEmail()); err != nil {
		errorsFields = append(errorsFields, fieldViolation("email", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	var newUser db.User

	// Выполняем создание пользователя в одной транзакции с записью задачи на верификацию почты в редиску
	err := db.MakeTx(ctx, func(q db.Querier) error {
		passwordHash, err := tools.GetPasswordHash(req.GetPassword())
		if err != nil {
			return status.Errorf(codes.Internal, "ошибка при хешировании: %s", err)
		}

		// Создание пользователя
		newUser, err = q.CreateUser(ctx, db.CreateUserParams{
			Name:         req.Name,
			PasswordHash: passwordHash,
			Email:        req.Email,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "не получилось создать пользователя: %s", err)
		}

		// Отдельно от создания пользователя создаём ещё и задачу (которую потом распределяем редиской)
		payload := &worker2.PayloadSendVerifyEmail{IdUser: newUser.IDUser}

		// Дополнительные конфигурации
		options := []asynq.Option{
			asynq.MaxRetry(4), // Максимальное количество повторений запроса при ошибках
			// (задержка нужна чтобы эта транзакция завершилась до того, как начнётся таска верификации почты)
			asynq.ProcessIn(10 * time.Second), // После какого времени процессору можно начать задачу
			asynq.Queue(worker2.QueueDefault), // Можем распределить важные задачи в отдельный поток (см. processor.go)
		}

		return server.taskDistributor.DistributeTaskVerifyEmail(ctx, payload, options...)
	})
	if err != nil {
		// Если пользователь уже есть, то выдаем ошибку
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			return nil, status.Errorf(codes.AlreadyExists, "пользователь с таким именем или email уже существует")
		}
		return nil, status.Errorf(codes.Internal, "не получилось создать пользователя: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convUser(newUser),
	}
	return response, nil
}

func validateGetUserRequest(req *pb.GetUserRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdUser())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_user", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Попытка получения данных из кеша
	cacheKey := fmt.Sprintf("user:%d", req.GetIdUser())
	var cachedResponse pb.GetUserResponse
	if err := server.cacher.GetCache(ctx, cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Проверка валидности запроса
	if err := validateGetUserRequest(req); err != nil {
		return nil, err
	}

	// Сам запрос
	user, err := server.querier.GetUser(ctx, req.GetIdUser())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить пользователя: %s", err)
	}

	response := &pb.GetUserResponse{
		User: convUser(user),
	}

	if err := server.cacher.SetCache(ctx, cacheKey, response); err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось сохранить данные в кеш: %s", err)
	}

	return response, nil
}

func validateGetManySortedUsersRequest(req *pb.GetManySortedUsersRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetPageNum())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_num", err))
	}

	if err := tools.ValidateNaturalNum(int(req.GetPageSize())); err != nil {
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
	users, err := server.querier.GetManySortedUsers(ctx, arg)
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
		if err := tools.ValidateString(req.GetName(), 1, 0); err != nil {
			errorsFields = append(errorsFields, fieldViolation("name", err))
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
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не авторизовался: %s", err)
	}

	// Чтобы изменить карму надо либо быть системой которая начисляем карму,
	// либо администратором (проверяем права)
	if req.Karma != nil && tools.HasPermission([]string{accessPayload.Role}, tools.UsualRole) {
		return nil, status.Errorf(codes.PermissionDenied, "у вас нет прав на изменение кармы")
	}

	arg := db.EditUserParams{
		IDUser: accessPayload.IDUser,
		Name:   req.GetName(),
		// Разделяем пустое значение и значение которое вообще не указывали
		// (Т.е. имеем возможность указать '' или 0 как валидный параметр (стереть значения), но не nil)
		Description: pgtype.Text{String: req.GetDescription(), Valid: req.Description != nil},
		Karma:       pgtype.Int4{Int32: req.GetKarma(), Valid: req.Karma != nil},
	}

	editedUser, err := server.querier.EditUser(ctx, arg)
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
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не авторизовался: %s", err)
	}

	deletedUser, err := server.querier.DeleteUser(ctx, payload.IDUser)
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

	if err := tools.ValidateString(req.GetPassword(), 1, 999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("password", err))
	}

	if err := tools.ValidateString(req.GetName(), 1, 99); err != nil {
		errorsFields = append(errorsFields, fieldViolation("name", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if err := validateLoginUserRequest(req); err != nil {
		return nil, err
	}

	user, err := server.querier.GetUserFromName(ctx, req.GetName())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "пользователь не найден")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить пользователя: %s", err)
	}

	if !tools.CheckPassword(req.GetPassword(), user.PasswordHash) {
		return nil, status.Errorf(codes.Unauthenticated, "неправильный пароль")
	}

	// Если с паролем со входом всё ок, то создаём новую сессию
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.IDUser,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать access token: %s", err)
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.IDUser,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать refresh token: %s", err)
	}

	info, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	_, err = server.querier.CreateSession(ctx, db.CreateSessionParams{
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

	// Удаляем просроченные сессии
	err = server.querier.DeleteExpiredSessions(ctx)
	if err != nil && err != sql.ErrNoRows {
		log.Err(err).Msg("сессии не удалились")
	}

	return response, nil
}

func validateRenewAccessTokenRequest(req *pb.RenewAccessTokenRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateString(req.GetRefreshToken(), 1, 999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("refresh_token", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	if err := validateRenewAccessTokenRequest(req); err != nil {
		return nil, err
	}

	// Только для авторизованных пользователей
	info, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	refreshPayload, errVerifyToken := server.tokenMaker.VerifyToken(req.GetRefreshToken())
	if errVerifyToken == token.ErrorInvalidToken {
		return nil, status.Errorf(codes.Internal, "недействительный токен")
	}

	oldSession, err := server.querier.DeleteSession(ctx, pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})
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
	if errVerifyToken != nil {
		return nil, status.Errorf(codes.Internal, errVerifyToken.Error())
	}

	newRefreshToken, newRefreshPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		refreshPayload.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать новый refresh token: %s", err)
	}
	newAccessToken, newAccessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось создать новый access token: %s", err)
	}

	_, err = server.querier.CreateSession(ctx, db.CreateSessionParams{
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

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdUser())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_user", err))
	}

	if err := tools.ValidateString(req.GetSecretKey(), 32, 0); err != nil {
		errorsFields = append(errorsFields, fieldViolation("secret_key", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if err := validateVerifyEmailRequest(req); err != nil {
		return nil, err
	}

	// ***пользователь перешёл по ссылке, верификация пройдена***

	err := db.MakeTx(ctx, func(queries db.Querier) error {
		verifyRequest, err := server.querier.GetVerifyRequest(ctx, req.GetIdUser())
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("запроса на верификацию почты нету")
			}
			return err
		}

		if verifyRequest.SecretKey != req.SecretKey {
			return fmt.Errorf("неправильный секретный ключ")
		}

		// Удаляем неактуальный запрос
		_, err = server.querier.DeleteVerifyRequest(ctx, req.GetIdUser())
		if err != nil {
			return err
		}

		// Если пользователь очень долго переходил по ссылке и она успела просрочиться
		if time.Now().After(verifyRequest.ExpiredAt.Time) {
			return status.Errorf(
				codes.Unauthenticated,
				"время для подтверждения почты истекло, пройдите снова процедуру верификации почты",
			)
		}

		_, err = server.querier.SetEmailIsVerified(ctx, verifyRequest.IDUser)
		if err != nil {
			return status.Errorf(codes.Internal, "%s", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Удаляем просроченные верификации
	err = server.querier.DeleteExiredRequests(ctx)
	if err != nil && err != sql.ErrNoRows {
		log.Err(err).Msg("запросы на верификацию почт не удалились")
	}

	return &pb.VerifyEmailResponse{}, nil
}
