package api

import (
	"context"
	"database/sql"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

// createUserRequest Поле Name обязательно (required), а Description заполним сами в createUser
// (json == что логично, берём данные из json'а в теле запроса)
type createUserRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (proc *Process) createUser(ctx *gin.Context) {
	var req createUserRequest

	// Проверяем чтобы все теги соответсвовали (в gin есть валидатор)
	// (в нашем случае чтобы было поле Name, иначе выдаём ошибку в JSON'е)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Если нету описания, то добавляем всё сами
	if req.Description == "" {
		req.Description = "Education is the cool site"
	}

	// Создаём пользователя
	arg := db.CreateUserParams{
		Name:        req.Name,
		Description: req.Description,
	}
	user, err := proc.queries.CreateUser(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// createUserRequest Нам нужен парамерт URI id_user который >= 1
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
	user, err := proc.queries.GetUser(context.Background(), req.IDUser)
	if err != nil {
		// Если у нас просто нет такого пользователя, то выдаём другую ошибку
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// getManyUsersRequest Сколько пользователей на страничке
// (form == берём данные из uri, которые идут после "?" (типа: /user?page_size=20&page_num=1))
type getManyUsersRequest struct {
	IDUser      bool  `json:"id_user"`
	Name        bool  `json:"name"`
	Description bool  `json:"description"`
	Karma       bool  `json:"karma"`
	PageSize    int32 `json:"page_size" binding:"required,min=1"`
	PageNum     int32 `json:"page_num" binding:"required,min=1"`
}

func (proc *Process) getManyUsers(ctx *gin.Context) {
	var req getManyUsersRequest

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
	users, err := proc.queries.GetManySortedUsers(context.Background(), arg)
	if err != nil {
		// Если у нас просто нет таких пользователей, то выдаём другую ошибку
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// Надо разделить данные которые получаем с url и данные которые получаем с uri
type editUserParamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Karma       int32  `json:"karma"`
	IDUser      int32  `json:"id_user" binding:"required,min=1"`
}

func (proc *Process) editUserParam(ctx *gin.Context) {
	var req editUserParamRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Изменяем параметр(ы) пользователя
	arg := db.EditUserParamParams{
		IDUser:      req.IDUser,
		Name:        req.Name,
		Description: req.Description,
		Karma:       req.Karma,
	}

	editedUser, err := proc.queries.EditUserParam(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, editedUser)
}

type deleteUserRequest struct {
	IDUser int32 `uri:"id_user" binding:"required"`
}

func (proc *Process) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	// Проверяем чтобы все теги соответствовали
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Удаляем пользователя
	deletedUser, err := proc.queries.DeleteUser(context.Background(), req.IDUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, deletedUser)
}
