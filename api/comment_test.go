package api

import (
	"context"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/my_mocks"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"reflect"
	"testing"
	"time"
)

func getRandomComment() db.Comment {
	return db.Comment{
		IDComment: tools.GetRandomUint(),
		Author:    tools.GetRandomUint(),
		Text:      tools.GetRandomString(1),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		EditedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

func TestGetComment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	comment := getRandomComment()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetComment(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(comment.IDComment)).
		Times(1).
		Return(comment, nil)

	server := &Server{querier: mockQueries}

	req := &pb.GetCommentRequest{IdComment: comment.IDComment}
	resp, err := server.GetComment(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, comment.Text, resp.GetComment().Text)
	require.Equal(t, comment.Author, resp.GetComment().Author)
	require.Equal(t, comment.IDComment, resp.GetComment().IdComment)
	require.Equal(t, comment.EditedAt.Time.Local(), resp.GetComment().EditedAt.AsTime().Local())
	require.Equal(t, comment.CreatedAt.Time.Local(), resp.GetComment().CreatedAt.AsTime().Local())
}

func TestGetComment_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	server := &Server{querier: mockQueries}

	req := &pb.GetCommentRequest{IdComment: -99_999_999}
	resp, err := server.GetComment(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
	require.Nil(t, resp)
}

// Тест на ненайденный коммент
func TestGetComment_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	mockQueries.EXPECT().
		GetComment(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(int32(99_999_999))).
		Times(1).
		Return(db.Comment{}, pgx.ErrNoRows)

	server := &Server{querier: mockQueries}

	req := &pb.GetCommentRequest{IdComment: 99_999_999}
	resp, err := server.GetComment(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, state.Code())
	require.Nil(t, resp)
}

// Тест на недопустимый идентификатор комментария
func TestGetComment_Validate2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	server := &Server{querier: mockQueries}

	req := &pb.GetCommentRequest{IdComment: -99_999_999}
	resp, err := server.GetComment(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
	require.Nil(t, resp)
}

// Тест на создание коммента
func TestCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockTokenMaker := my_mocks.NewMockMaker(ctrl)

	comm := getRandomComment()

	mockQueries.EXPECT().
		CreateComment(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Any()).
		Times(1).
		Return(comm, nil)

	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Any()).
		Times(1).
		Return(&token.Payload{IDUser: tools.GetRandomUint()}, nil)

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: supportedAuthType + " abc"}),
	)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}

	req := &pb.CreateCommentRequest{
		Text:      tools.GetRandomString(),
		IdArticle: tools.GetRandomUint(),
	}
	resp, err := server.CreateComment(loginCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, comm.Text, resp.GetComment().GetText())
}

// Тест на удаление комментария
func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockTokenMaker := my_mocks.NewMockMaker(ctrl)

	idUser := tools.GetRandomUint()
	idComment := tools.GetRandomUint()

	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Any()).
		Times(1).
		Return(&token.Payload{IDUser: idUser}, nil)

	mockQueries.EXPECT().
		DeleteComment(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(idComment)).
		Times(1).
		Return(db.Article{}, nil)
	mockQueries.EXPECT().
		GetComment(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(idComment)).
		Times(1).
		Return(db.Comment{Author: idUser}, nil)

	loginCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{authHeader: supportedAuthType + " abc"}),
	)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}

	req := &pb.DeleteCommentRequest{IdComment: idComment}
	resp, err := server.DeleteComment(loginCtx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, req.GetIdComment(), idComment)
}

// Тест на недопустимый идентификатор комментария
func TestDeleteComment_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	server := &Server{querier: mockQueries}

	req := &pb.DeleteCommentRequest{IdComment: -99_999_999}
	resp, err := server.DeleteComment(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
	require.Nil(t, resp)
}

// Тест на недопустимый идентификатор комментария
func TestEditComment_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	server := &Server{querier: mockQueries}

	txt := tools.GetRandomString(1)
	req := &pb.EditCommentRequest{IdComment: -99_999_999, Text: &txt}
	resp, err := server.EditComment(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
	require.Nil(t, resp)
}

// Тест на недопустимый идентификатор комментария
func TestCreateComment_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueries := my_mocks.NewMockQuerier(ctrl)

	server := &Server{querier: mockQueries}

	req := &pb.CreateCommentRequest{Text: tools.GetRandomString(1)}
	resp, err := server.CreateComment(context.Background(), req)
	require.Error(t, err)
	state, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, state.Code())
	require.Nil(t, resp)
}
