package api_gin

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// createUserRequest Поле Name обязательно (required) и без спецсимволов (alphanum),
// а также минимальная длина пароля - 6 символов
// (json == что логично, берём данные из json'а в теле запроса)
type createUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	// Проверяем чтобы все теги соответсвовали (в gin есть валидатор)
	// (в нашем случае чтобы было поле Username, иначе выдаём ошибку в JSON'е)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordHash, err := tools.GetPasswordHash(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Создаём пользователя
	arg := db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	user, err := server.querier.CreateUser(ctx, arg)
	if err != nil {
		// Если пользователь с таким именем уже есть, то выдаем ошибку
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			ctx.JSON(http.StatusConflict, errorResponse(errors.New("пользователь с таким именем или email уже существует")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userToResponse(&user)
	ctx.JSON(http.StatusOK, user)
}

// getUserRequest Нам нужен парамерт URI id_user который >= 1
// (uri == берём данные из uri (типа: user/42))
type getUserRequest struct {
	IDUser int32 `uri:"id_user" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	// Проверяем чтобы все теги соответсвовали (в gin есть валидатор)
	// (в нашем случае чтобы было поле IDUser, иначе выдаём ошибку в JSON'е)
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем пользователя
	user, err := server.querier.GetUser(ctx, req.IDUser)
	if err != nil {
		// Если у нас просто нет такого пользователя, то выдаём другую ошибку
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userToResponse(&user)
	ctx.JSON(http.StatusOK, user)
}

// getManySortedUsersRequest Сколько а как отсортированных пользователей на страничке
// (form == берём данные из uri, которые идут после "?" (типа: /user?page_size=20&page_num=1))
type getManySortedUsersRequest struct {
	IDUser      bool  `json:"id_user"`
	Name        bool  `json:"name"`
	Description bool  `json:"description"`
	Karma       bool  `json:"karma"`
	PageSize    int32 `json:"page_size" binding:"required,min=1"`
	PageNum     int32 `json:"page_num" binding:"required,min=1"`
}

func (server *Server) getManySortedUsers(ctx *gin.Context) {
	var req getManySortedUsersRequest

	// Проверяем чтобы все теги соответсвовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем пользователей
	arg := db.GetManySortedUsersParams{
		IDUser:      req.IDUser,
		Name:        req.Name,
		Description: req.Description,
		Karma:       req.Karma,
		Limit:       req.PageSize,
		Offset:      (req.PageNum - 1) * req.PageSize,
	}
	users, err := server.querier.GetManySortedUsers(ctx, arg)
	if err != nil {
		// Если у нас просто нет таких пользователей, то выдаём другую ошибку
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Превращаем в нужный формат
	for i, u := range users {
		users[i] = *userToResponse(&u)
	}

	ctx.JSON(http.StatusOK, users)
}

// Надо разделить данные которые получаем с url и данные которые получаем с uri
// (оно разделяет пустые строки, полученные из json'а от пустых строк вставленные go автоматически
// (поэтому тут такие костыли))
type editUserRequest struct {
	Name        string
	Description pgtype.Text
	Karma       pgtype.Int4
}

func (server *Server) editUser(ctx *gin.Context) {
	var body map[string]interface{}

	// Связываем все поля которые нам нужны
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Делаем операцию только для авторизованного пользователя
	t, exist := ctx.Get(authPayloadKey)
	if !exist {
		t = ctx.Request.Header.Get(authPayloadKey)
	}
	tkn := t.(string)
	payload, err := server.tokenMaker.VerifyToken(tkn)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Тут мы разделяем поля которые содержат пустые строки от полей которые вообще не указаны в теле
	var req editUserRequest
	bodyText := make(map[string]string)
	bodyInt := make(map[string]int32)
	for key, val := range body {
		switch v := val.(type) {
		case string:
			bodyText[key] = v
		case float64:
			bodyInt[key] = int32(v)
		}
	}

	if val, ok := bodyText["name"]; ok {
		req.Name = val
	}

	val, ok := bodyText["description"]
	req.Description = pgtype.Text{String: val, Valid: ok}

	num, ok := bodyInt["karma"]
	req.Karma = pgtype.Int4{Int32: num, Valid: ok}

	// Изменяем параметр(ы) пользователя
	arg := db.EditUserParams{
		IDUser:      payload.IDUser,
		Name:        req.Name,
		Description: req.Description,
		Karma:       req.Karma,
	}

	editedUser, err := server.querier.EditUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userToResponse(&editedUser)
	ctx.JSON(http.StatusOK, editedUser)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	// Удаляем авторизованного пользователя
	t, exist := ctx.Get(authPayloadKey)
	if !exist {
		t = ctx.Request.Header.Get(authPayloadKey)
	}
	tkn := t.(string)
	payload, err := server.tokenMaker.VerifyToken(tkn)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	deletedUser, err := server.querier.DeleteUser(ctx, payload.IDUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userToResponse(&deletedUser)
	ctx.JSON(http.StatusOK, deletedUser)
}

// loginUserRequest Логиним пользователя
type loginUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

// loginUserResponse Отправляем токен
type loginUserResponse struct {
	User                  db.User   `json:"user"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshToken          string    `json:"resfresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}

// userToResponse Заменяем PasswordHash, т.к. передавать его не безопасно
func userToResponse(user *db.User) *db.User {
	user.PasswordHash = "wat u looking at :)"
	return user
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Входим в систему
	user, err := server.querier.GetUserFromName(ctx, req.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Проверяем пароль
	if !tools.CheckPassword(req.Password, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("неправильный пароль")))
		return
	}

	// Залогиненному пользователю даём access и refresh токены
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.IDUser,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.IDUser,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.querier.CreateSession(ctx, db.CreateSessionParams{
		IDSession:    pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true},
		IDUser:       user.IDUser,
		RefreshToken: refreshToken,
		ExpiredAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiredAt, Valid: true},
		ClientIp:     ctx.ClientIP(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshPayload.ExpiredAt,
		User:                  user,
	}

	userToResponse(&response.User)
	ctx.JSON(http.StatusOK, response)
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshToken          string    `json:"resfresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}

// renewAccessToken Входим по refresh токену и обновляем ОБА токена
func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, errVerifyToken := server.tokenMaker.VerifyToken(req.RefreshToken)

	// Пересоздаём и сессию (refresh токен) и access токен
	oldSession, err := server.querier.DeleteSession(ctx, pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})

	// Тут проверяем что этот refresh токен валидный
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else if oldSession.Blocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("сессия заблокирована")))
		return
	} else if oldSession.IDUser != refreshPayload.IDUser {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("некорректная сессия пользователя")))
		return
	} else if oldSession.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("некорректный refresh токен")))
		return
	}
	// Если токен просрочен или ip не соответствует ранее залогиневшему устройству, то перенаправляем на ввод пароля
	if errVerifyToken == token.ErrorExpiredToken || oldSession.ClientIp != ctx.ClientIP() {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("необходимо залогиниться"))
		return
	}

	// Создаём новые access и refresh токены (и сессию)
	newRefreshToken, newRefreshPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		refreshPayload.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	newAccessToken, newAccessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.querier.CreateSession(ctx, db.CreateSessionParams{
		IDSession:    pgtype.UUID{Bytes: newRefreshPayload.IDSession, Valid: true},
		IDUser:       newRefreshPayload.IDUser,
		RefreshToken: newRefreshToken,
		ExpiredAt:    pgtype.Timestamptz{Time: newRefreshPayload.ExpiredAt, Valid: true},
		ClientIp:     ctx.ClientIP(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:           newAccessToken,
		AccessTokenExpiredAt:  newAccessPayload.ExpiredAt,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiredAt: newRefreshPayload.ExpiredAt,
	})

	// Удаляем просроченные сессии
	err = server.querier.DeleteExpiredSessions(ctx)
	if err != nil || err != sql.ErrNoRows {
		log.Err(err).Msg("сессии не удалились")
	}
}
