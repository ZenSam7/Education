package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"
)

// createUserRequest Поле Name обязательно (required) и без спецсимволов (alphanum),
// а также минимальная длина пароля - 6 символов
// (json == что логично, берём данные из json'а в теле запроса)
type createUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (proc *Process) createUser(ctx *gin.Context) {
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
	user, err := proc.queries.CreateUser(ctx, arg)
	if err != nil {
		// Если пользователь с таким именем уже есть, то выдаем ошибку
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			ctx.JSON(http.StatusConflict, errorResponse(errors.New("user with this name or email already exists")))
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

func (proc *Process) getUser(ctx *gin.Context) {
	var req getUserRequest

	// Проверяем чтобы все теги соответсвовали (в gin есть валидатор)
	// (в нашем случае чтобы было поле IDUser, иначе выдаём ошибку в JSON'е)
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем пользователя
	user, err := proc.queries.GetUser(ctx, req.IDUser)
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

func (proc *Process) getManySortedUsers(ctx *gin.Context) {
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
	users, err := proc.queries.GetManySortedUsers(ctx, arg)
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
type editUserParamRequest struct {
	Description string `json:"description"`
	Karma       int32  `json:"karma"`
	Name        string `json:"name"`
}

func (proc *Process) editUserParam(ctx *gin.Context) {
	var req editUserParamRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Делаем операцию только для авторизованного пользователя
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)

	// Изменяем параметр(ы) пользователя
	arg := db.EditUserParams{
		IDUser:      payload.IDUser,
		Name:        req.Name,
		Description: req.Description,
		Karma:       req.Karma,
	}

	editedUser, err := proc.queries.EditUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userToResponse(&editedUser)
	ctx.JSON(http.StatusOK, editedUser)
}

func (proc *Process) deleteUser(ctx *gin.Context) {
	// Удаляем авторизованного пользователя
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)
	deletedUser, err := proc.queries.DeleteUser(ctx, payload.IDUser)
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
	Password string `json:"password" binding:"required,min=6"`
}

// loginUserResponse Отправляем токен
type loginUserResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshToken          string    `json:"resfresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
	User                  db.User   `json:"user"`
}

// userToResponse Заменяем PasswordHash, т.к. передавать его не безопасно
func userToResponse(user *db.User) *db.User {
	user.PasswordHash = "wat u looking at :)"
	return user
}

func (proc *Process) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Входим в систему
	user, err := proc.queries.GetUserForName(ctx, req.Name)
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
	accessToken, accessTokenPayload, err := proc.tokenMaker.CreateToken(user.IDUser, proc.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := proc.tokenMaker.CreateToken(user.IDUser, proc.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = proc.queries.CreateSession(ctx, db.CreateSessionParams{
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
func (proc *Process) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, errVerifyToken := proc.tokenMaker.VerifyToken(req.RefreshToken)

	// Пересоздаём и сессию (refresh токен) и access токен
	oldSession, err := proc.queries.DeleteSession(ctx, pgtype.UUID{Bytes: refreshPayload.IDSession, Valid: true})

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
	newRefreshToken, newRefreshPayload, err := proc.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		proc.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	newAccessToken, newAccessPayload, err := proc.tokenMaker.CreateToken(
		refreshPayload.IDUser,
		proc.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = proc.queries.CreateSession(ctx, db.CreateSessionParams{
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
}
