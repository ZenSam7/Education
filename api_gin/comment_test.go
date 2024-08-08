package api_gin

import (
	"context"
	"database/sql"
	"encoding/json"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/my_mocks"
	"github.com/ZenSam7/Education/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func getRandomComment() db.Comment {
	return db.Comment{
		IDComment: tools.GetRandomUint(),
		Author:    tools.GetRandomUint(),
		Text:      tools.GetRandomString(1),
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
	router := gin.Default()
	router.GET("/comment/:id_comment", server.getComment)

	req, _ := http.NewRequest("GET", "/comment/"+strconv.Itoa(int(comment.IDComment)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var gotComment db.Comment
	err := json.Unmarshal(w.Body.Bytes(), &gotComment)
	require.NoError(t, err)
	require.Equal(t, comment.Text, gotComment.Text)
	require.Equal(t, comment.IDComment, gotComment.IDComment)
	require.Equal(t, comment.Author, gotComment.Author)
}

func TestGetComment_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idComment := int32(99999999)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetComment(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(idComment)).
		Times(1).
		Return(db.Comment{}, sql.ErrNoRows)

	server := &Server{querier: mockQueries}
	router := gin.Default()
	router.GET("/comment/:id_comment", server.getComment)

	req, _ := http.NewRequest("GET", "/comment/"+strconv.Itoa(int(idComment)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetComment_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := &Server{}
	router := gin.Default()
	router.GET("/comment/:id_comment", server.getComment)

	req, _ := http.NewRequest("GET", "/comment/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
