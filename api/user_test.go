package api

import (
	"context"
	"database/sql"
	"github.com/ZenSam7/Education/db/mockdb"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/tools"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

// createRandomUser Создаём случайного пользователя и возвращаем функцию для закрытия соединения
func createRandomUser() (db.User, *db.Queries, func()) {
	queries, closeConn := db.MakeQueries()

	arg := db.CreateUserParams{
		Name:         tools.GetRandomString(),
		Email:        tools.GetRandomEmail(),
		PasswordHash: tools.GetRandomHash(),
	}
	newUser, _ := queries.CreateUser(context.Background(), arg)

	return newUser, queries, closeConn
}

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

	mockQueries := mockdb.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUser(gomock.Any(), gomock.Eq(user.IDUser)).
		Times(1).
		Return(user, nil)

	server := &Server{
		querier: (db.Querier)(mockQueries),
	}

	req := &pb.GetUserRequest{IdUser: user.IDUser}
	resp, err := server.GetUser(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.Name, resp.GetUser().Name)
	assert.Equal(t, user.Email, resp.GetUser().Email)
	assert.Equal(t, user.Karma, resp.GetUser().Karma)
	assert.Equal(t, user.IDUser, resp.GetUser().IdUser)
	assert.Equal(t, user.Description.String, resp.GetUser().Description)
}

// Тест на недопустимый идентификатор пользователя
func TestGetUser_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := mockdb.NewMockQuerier(ctrl)

	server := &Server{
		querier: (db.Querier)(mockQueries),
	}

	req := &pb.GetUserRequest{IdUser: -999999}
	resp, err := server.GetUser(context.Background(), req)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, resp)
}

// Тест на ненайденного пользователя
func TestGetUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idUser := int32(1)

	mockQueries := mockdb.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUser(gomock.Any(), idUser).
		Times(1).
		Return(getRandomUser(), sql.ErrNoRows)

	server := &Server{
		querier: (db.Querier)(mockQueries),
	}

	req := &pb.GetUserRequest{IdUser: idUser}
	resp, err := server.GetUser(context.Background(), req)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, resp)
}
