package api_gin

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/my_mocks"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
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

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := getRandomUser()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(user.IDUser)).
		Times(1).
		Return(user, nil)

	server := &Server{querier: mockQueries}
	router := gin.Default()
	router.GET("/user/:id_user", server.getUser)

	req, _ := http.NewRequest("GET", "/user/"+strconv.Itoa(int(user.IDUser)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var gotUser db.User
	err := json.Unmarshal(w.Body.Bytes(), &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Karma, gotUser.Karma)
	require.Equal(t, user.IDUser, gotUser.IDUser)
	require.Equal(t, user.Description.String, gotUser.Description.String)
}

func TestGetUser_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := &Server{querier: my_mocks.NewMockQuerier(ctrl)}
	router := gin.Default()
	router.GET("/user/:id_user", server.getUser)

	req, _ := http.NewRequest("GET", "/user/-99999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idUser := int32(99999999)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUser(gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()), gomock.Eq(idUser)).
		Times(1).
		Return(db.User{}, sql.ErrNoRows)

	server := &Server{querier: mockQueries}
	router := gin.Default()
	router.GET("/user/:id_user", server.getUser)

	req, _ := http.NewRequest("GET", "/user/"+strconv.Itoa(int(idUser)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateUser_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{
		router: gin.Default(),
	}

	server.router.POST("/users", server.createUser)

	recorder := httptest.NewRecorder()
	reqBody := createUserRequest{
		Name:     "",              // Неверное значение для имени
		Email:    "invalid-email", // Неверное значение для email
		Password: "",              // Неверное значение для пароля
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetManySortedUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetManySortedUsers(gomock.Any(), db.GetManySortedUsersParams{
			Offset: 0, // PageNum - 1 * PageSize
			Limit:  10,
		}).
		Times(1).
		Return([]db.User{getRandomUser()}, nil)

	server := &Server{
		router:  gin.Default(),
		querier: mockQueries,
	}

	server.router.POST("/users", server.getManySortedUsers)

	recorder := httptest.NewRecorder()
	reqBody := getManySortedUsersRequest{
		PageNum:  1,
		PageSize: 10,
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetManySortedUsers_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := &Server{
		router: gin.Default(),
	}

	server.router.POST("/users", server.getManySortedUsers)

	recorder := httptest.NewRecorder()
	reqBody := getManySortedUsersRequest{
		PageNum:  -1,  // Неверное значение для номера страницы
		PageSize: -10, // Неверное значение для размера страницы
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestEditUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := getRandomUser()
	editedUser := getRandomUser()
	editedUser.Name = "newname"
	editedUser.Description = pgtype.Text{String: "newdescription", Valid: true}
	editedUser.Karma = tools.GetRandomInt()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		EditUser(gomock.Any(), db.EditUserParams{
			IDUser:      user.IDUser,
			Name:        "newname",
			Description: pgtype.Text{String: "newdescription", Valid: true},
			Karma:       pgtype.Int4{Int32: editedUser.Karma, Valid: true},
		}).
		Times(1).
		Return(editedUser, nil)
	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Any()).
		Times(1).
		Return(&token.Payload{IDUser: user.IDUser}, nil)

	// Создаем новый Gin роутер
	r := gin.Default()
	server := &Server{
		router:     r,
		querier:    mockQueries,
		tokenMaker: mockTokenMaker,
	}

	// Устанавливаем маршрут для редактирования пользователя
	r.PATCH("/users", server.editUser)

	// Создаем тестовый запрос
	recorder := httptest.NewRecorder()
	reqBody := map[string]interface{}{
		"name":        "newname",
		"description": "newdescription",
		"karma":       editedUser.Karma,
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPatch, "/users", bytes.NewReader(body))
	require.NoError(t, err)

	// Обрабатываем запрос
	r.ServeHTTP(recorder, request)

	// Проверяем результат
	require.Equal(t, http.StatusOK, recorder.Code)

	var gotUser db.User
	err = json.Unmarshal(recorder.Body.Bytes(), &gotUser)
	require.NoError(t, err)

	require.Equal(t, editedUser.Name, gotUser.Name)
	require.Equal(t, editedUser.Karma, gotUser.Karma)
	require.Equal(t, editedUser.Description.String, gotUser.Description.String)
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := getRandomUser()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		DeleteUser(gomock.Any(), gomock.Eq(user.IDUser)).
		Times(1).
		Return(user, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(gomock.Any()).
		Times(1).
		Return(&token.Payload{IDUser: user.IDUser}, nil)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	router := gin.Default()
	router.DELETE("/user", server.deleteUser)

	req, _ := http.NewRequest(http.MethodDelete, "/user", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var responseUser db.User
	err := json.NewDecoder(resp.Body).Decode(&responseUser)
	require.NoError(t, err)
	require.Equal(t, user.IDUser, responseUser.IDUser)
}

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := loginUserRequest{
		Name:     "testname",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	user := getRandomUser()
	user.PasswordHash, _ = tools.GetPasswordHash(reqBody.Password)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUserFromName(gomock.Any(), gomock.Eq(reqBody.Name)).
		Times(1).
		Return(user, nil)
	mockQueries.EXPECT().
		CreateSession(gomock.Any(), gomock.Any()).
		Times(1).
		Return(db.Session{}, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		CreateToken(user.IDUser, user.Role, gomock.Any()).
		Times(2).
		Return("token", &token.Payload{IDUser: user.IDUser}, nil)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	router := gin.Default()
	router.POST("/login", server.loginUser)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var response loginUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, user.IDUser, response.User.IDUser)
}

func TestLoginUser_ValidationError(t *testing.T) {
	reqBody := loginUserRequest{
		Name:     "",
		Password: "",
	}
	body, _ := json.Marshal(reqBody)

	server := &Server{}
	router := gin.Default()
	router.POST("/login", server.loginUser)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := loginUserRequest{
		Name:     "nonexistentuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUserFromName(gomock.Any(), gomock.Eq(reqBody.Name)).
		Times(1).
		Return(db.User{}, sql.ErrNoRows)

	server := &Server{querier: mockQueries}
	router := gin.Default()
	router.POST("/login", server.loginUser)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := loginUserRequest{
		Name:     "testname",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	user := getRandomUser()

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		GetUserFromName(gomock.Any(), gomock.Eq(reqBody.Name)).
		Times(1).
		Return(user, nil)

	server := &Server{querier: mockQueries}
	router := gin.Default()
	router.POST("/login", server.loginUser)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestRenewAccessToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := renewAccessTokenRequest{
		RefreshToken: "InvalidToken",
	}
	body, _ := json.Marshal(reqBody)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(reqBody.RefreshToken).
		Times(1).
		Return(nil, token.ErrorInvalidToken)

	server := &Server{tokenMaker: mockTokenMaker}
	router := gin.Default()
	router.POST("/renew", server.renewAccessToken)

	req, _ := http.NewRequest(http.MethodPost, "/renew", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestRenewAccessToken_SessionBlocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := renewAccessTokenRequest{
		RefreshToken: "validRefreshToken",
	}
	body, _ := json.Marshal(reqBody)

	user := getRandomUser()
	refreshPayload := &token.Payload{
		IDUser:    user.IDUser,
		IDSession: [16]byte(uuid.New()),
	}
	blockedSession := db.Session{
		IDUser:       user.IDUser,
		RefreshToken: reqBody.RefreshToken,
		Blocked:      true,
		ClientIp:     "127.0.0.1",
	}

	mockQueries := my_mocks.NewMockQuerier(ctrl)
	mockQueries.EXPECT().
		DeleteSession(gomock.Any(), gomock.Eq(pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})).
		Times(1).
		Return(blockedSession, nil)

	mockTokenMaker := my_mocks.NewMockMaker(ctrl)
	mockTokenMaker.EXPECT().
		VerifyToken(reqBody.RefreshToken).
		Times(1).
		Return(refreshPayload, nil)

	server := &Server{querier: mockQueries, tokenMaker: mockTokenMaker}
	router := gin.Default()
	router.POST("/renew", server.renewAccessToken)

	req, _ := http.NewRequest(http.MethodPost, "/renew", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)
}
