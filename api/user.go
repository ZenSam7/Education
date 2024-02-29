package api

import (
	"context"
	"database/sql"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

// createUserRequest Поле Name обязательно, а Description заполним сами в createUser
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

	//// Если нету описания, то добавляем всё сами
	//if req.Description == "" {
	//	req.Description = "Education is the cool site"
	//}

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
type getManyUsersRequest struct {
	IDUser      bool  `form:"id_user"`
	Name        bool  `form:"name"`
	Description bool  `form:"description"`
	Karma       bool  `form:"karma"`
	PageSize    int32 `form:"page_size" binding:"required,min=1"`
	PageNum     int32 `form:"page_num" binding:"required,min=1"`
}

func (proc *Process) getManyUsers(ctx *gin.Context) {
	var req getManyUsersRequest

	// Проверяем чтобы все теги соответсвовали
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем пользователя
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
