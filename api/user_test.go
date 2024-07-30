package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/my_mocks"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func getRandomUser() db.User {
	return db.User{
		IDUser:       tools.GetRandomUint(),
		Name:         tools.GetRandomString(1),
		Email:        tools.GetRandomEmail(),
		PasswordHash: tools.GetRandomHash(),
		Description:  pgtype.Text{String: tools.GetRandomString(), Valid: true},
		Karma:        tools.GetRandomInt(),
	}
}

// Тест на успешный запрос
func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := getRandomUser()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUser(gomock.Any(), gomock.Eq(user.IDUser)).
		Times(1).
		Return(user, nil)

	server := &Server{querier: mockQueries}

	req := &pb.GetUserRequest{IdUser: user.IDUser}
	resp, err := server.GetUser(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, user.Name, resp.GetUser().Name)
	require.Equal(t, user.Email, resp.GetUser().Email)
	require.Equal(t, user.Karma, resp.GetUser().Karma)
	require.Equal(t, user.IDUser, resp.GetUser().IdUser)
	require.Equal(t, user.Description.String, resp.GetUser().Description)
}

// Тест на недопустимый идентификатор пользователя
func TestGetUser_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	server := &Server{querier: mockQueries}

	req := &pb.GetUserRequest{IdUser: -99_999_999}
	resp, err := server.GetUser(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
	require.Nil(t, resp)
}

// Тест на ненайденного пользователя
func TestGetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idUser := int32(99_999_999)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUser(gomock.Any(), gomock.Eq(idUser)).
		Times(1).
		Return(getRandomUser(), sql.ErrNoRows)

	server := &Server{querier: mockQueries}

	req := &pb.GetUserRequest{IdUser: idUser}
	resp, err := server.GetUser(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, state.Code())
	require.Nil(t, resp)
}

// Тест на ошибку валидации при создании пользователя
func TestCreateUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.CreateUserRequest{
		Name:     "",              // Неверное значение для имени
		Email:    "invalid-email", // Неверное значение для email
		Password: "",              // Неверное значение для пароля
	}

	server := &Server{}
	resp, err := server.CreateUser(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
}

// Тест на успешное получение пользователей
func TestGetManySortedUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.GetManySortedUsersRequest{
		PageNum:  1,
		PageSize: 10,
	}

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetManySortedUsers(gomock.Any(), gomock.Eq(db.GetManySortedUsersParams{
			Offset: (req.PageNum - 1) * req.PageSize, Limit: req.GetPageSize(),
		})).
		Times(1).
		Return([]db.User{getRandomUser()}, nil)

	server := &Server{querier: mockQueries}
	resp, err := server.GetManySortedUsers(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetUsers(), 1)
}

// Тест на ошибку валидации при получении пользователей
func TestGetManySortedUsers_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.GetManySortedUsersRequest{
		PageNum:  -1,  // Неверное значение для номера страницы
		PageSize: -10, // Неверное значение для размера страницы
	}

	server := &Server{}
	resp, err := server.GetManySortedUsers(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
}

// Тест на успешное редактирование пользователя
func TestEditUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.EditUserRequest{
		Name:        proto.String("newname"),
		Description: proto.String("newdescription"),
		Karma:       proto.Int32(tools.GetRandomInt()),
	}

	user := getRandomUser()
	editedUser := getRandomUser()
	editedUser.Name = req.GetName()
	editedUser.Description = pgtype.Text{String: req.GetDescription(), Valid: true}
	editedUser.Karma = req.GetKarma()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		EditUser(gomock.Any(), db.EditUserParams{
			IDUser: user.IDUser, Name: req.GetName(),
			Description: pgtype.Text{String: req.GetDescription(), Valid: true},
			Karma:       pgtype.Int4{Int32: req.GetKarma(), Valid: true},
		}).
		Times(1).
		Return(user, nil)

	// Чтобы не логиниться, мокаем tokenMaker и что-то кладём в качестве токена
	// авторизации чтобы metadata работала
	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Any()).
		Times(1).
		Return(&token.Payload{IDUser: user.IDUser}, nil)

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: fmt.Sprintf("%s abc", supportedAuthType)}),
	)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}

	resp, err := server.EditUser(loginCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req.GetName(), editedUser.Name)
	require.Equal(t, req.GetKarma(), editedUser.Karma)
	require.Equal(t, req.GetDescription(), editedUser.Description.String)
}

// Тест на ошибку валидации при редактировании пользователя
func TestEditUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.EditUserRequest{
		Name: proto.String(""), // Неверное значение для имени
	}

	server := &Server{}
	resp, err := server.EditUser(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := getRandomUser()

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: fmt.Sprintf("%s abc", supportedAuthType)}),
	)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		DeleteUser(gomock.Eq(loginCtx), gomock.Eq(user.IDUser)).
		Times(1).
		Return(user, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Any()).
		Times(1).
		Return(&token.Payload{IDUser: user.IDUser}, nil)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	req := &pb.DeleteUserRequest{}
	resp, err := server.DeleteUser(loginCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, user.IDUser, resp.GetUser().IdUser)
}

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.LoginUserRequest{
		Name:     "testname",
		Password: "password123",
	}

	user := getRandomUser()
	user.PasswordHash, _ = tools.GetPasswordHash(req.GetPassword())

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: fmt.Sprintf("%s abc", supportedAuthType)}),
	)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUserFromName(loginCtx, gomock.Eq(req.GetName())).
		Times(1).
		Return(user, nil)
	mockQueries.EXPECT().
		CreateSession(loginCtx, gomock.Any()).
		Times(1).
		Return(db.Session{}, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		CreateToken(gomock.Any(), gomock.Any()).
		Times(2).
		Return("token", &token.Payload{IDUser: user.IDUser}, nil)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	resp, err := server.LoginUser(loginCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, user.IDUser, resp.GetUser().IdUser)
}

// Тест на ошибку валидации при входе пользователя
func TestLoginUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.LoginUserRequest{
		Name:     "",
		Password: "",
	}

	server := &Server{}
	resp, err := server.LoginUser(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
}

// Тест на ошибку: пользователь не найден
func TestLoginUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.LoginUserRequest{
		Name:     "nonexistentuser",
		Password: "password123",
	}

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUserFromName(gomock.Any(), gomock.Eq(req.GetName())).
		Times(1).
		Return(db.User{}, sql.ErrNoRows)

	server := &Server{querier: mockQueries}
	resp, err := server.LoginUser(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

// Тест на ошибку: неправильный пароль
func TestLoginUser_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.LoginUserRequest{
		Name:     "testname",
		Password: "wrongpassword",
	}

	user := getRandomUser()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUserFromName(gomock.Any(), gomock.Eq(req.GetName())).
		Times(1).
		Return(user, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	resp, err := server.LoginUser(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestRenewAccessToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.RenewAccessTokenRequest{
		RefreshToken: "validRefreshToken",
	}

	user := getRandomUser()
	refreshPayload := &token.Payload{
		IDUser:    user.IDUser,
		IDSession: [16]byte(uuid.New()),
	}
	oldSession := db.Session{
		IDUser:       user.IDUser,
		RefreshToken: req.GetRefreshToken(),
		ClientIp:     "1234.1234.1234.1234",
	}

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{
			authHeader:          fmt.Sprintf("%s abc", supportedAuthType),
			xForwardedForHeader: oldSession.ClientIp,
		}),
	)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		DeleteSession(loginCtx, gomock.Eq(pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})).
		Times(1).
		Return(oldSession, nil)
	mockQueries.EXPECT().
		CreateSession(loginCtx, gomock.Any()).
		Times(1).
		Return(db.Session{}, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(req.GetRefreshToken()).
		Times(1).
		Return(refreshPayload, nil)
	mockTokenMaker.EXPECT().
		CreateToken(gomock.Any(), gomock.Any()).
		Times(2).
		Return("newToken", refreshPayload, nil)

	server := &Server{
		querier:    mockQueries,
		tokenMaker: mockTokenMaker,
		config:     tools.Config{RefreshTokenDuration: time.Hour, AccessTokenDuration: time.Minute},
	}
	resp, err := server.RenewAccessToken(loginCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "newToken", resp.GetAccessToken())
}

func TestRenewAccessToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.RenewAccessTokenRequest{
		RefreshToken: "InvalidToken",
	}

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{
			authHeader: fmt.Sprintf("%s abc", supportedAuthType),
		}),
	)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Eq(req.RefreshToken)).
		Times(1).
		Return(nil, token.ErrorInvalidToken)

	server := &Server{tokenMaker: mockTokenMaker}
	resp, err := server.RenewAccessToken(loginCtx, req)
	require.Error(t, err)
	require.Nil(t, resp)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, state.Code())
}

func TestRenewAccessToken_SessionBlocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.RenewAccessTokenRequest{
		RefreshToken: "validRefreshToken",
	}

	user := getRandomUser()
	refreshPayload := &token.Payload{
		IDUser:    user.IDUser,
		IDSession: [16]byte(uuid.New()),
	}
	blockedSession := db.Session{
		IDUser:       user.IDUser,
		RefreshToken: req.GetRefreshToken(),
		Blocked:      true,
		ClientIp:     "127.0.0.1",
	}

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: fmt.Sprintf("%s abc", supportedAuthType)}),
	)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		DeleteSession(loginCtx, gomock.Eq(pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})).
		Times(1).
		Return(blockedSession, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(req.GetRefreshToken()).
		Times(1).
		Return(refreshPayload, nil)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	resp, err := server.RenewAccessToken(loginCtx, req)
	require.Error(t, err)
	require.Nil(t, resp)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, state.Code())
}
