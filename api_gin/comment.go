package api_gin

import (
	"context"
	"errors"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getCommentsOfArticleRequest struct {
	IDArticle int32 `uri:"id_article" binding:"required,min=1"`
	PageNum   int32 `uri:"page_num" binding:"required,min=1"`
	PageSize  int32 `uri:"page_size" binding:"required,min=1"`
}

func (server *Server) getCommentsOfArticle(ctx *gin.Context) {
	var req getCommentsOfArticleRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Возвращаем комментарии
	args := db.GetCommentsOfArticleParams{
		IDArticle: req.IDArticle,
		Limit:     req.PageSize,
		Offset:    (req.PageNum - 1) * req.PageSize,
	}
	comments, err := server.querier.GetCommentsOfArticle(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

type createCommentRequest struct {
	IDArticle int32  `json:"id_article" binding:"required,min=1"`
	Text      string `json:"text" binding:"required"`
	Author    int32  `json:"author" binding:"required"`
}

func (server *Server) createComment(ctx *gin.Context) {
	var req createCommentRequest

	// Комментарии может создать только авторизованный пользователь
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Создаем комментарий
	args := db.CreateCommentParams{
		IDArticle: req.IDArticle,
		Text:      req.Text,
		Author:    payload.IDUser,
	}
	comment, err := server.querier.CreateComment(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

type deleteCommentRequest struct {
	IDComment int32 `uri:"id_comment" binding:"required,min=1"`
}

func isAuthorComment(IDUser, IDComment int32, server *Server) bool {
	comment, _ := server.querier.GetComment(context.Background(), IDComment)
	return comment.Author == IDUser
}

func (server *Server) deleteComment(ctx *gin.Context) {
	var req deleteCommentRequest

	// Выявляем автора
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if !isAuthorComment(payload.IDUser, req.IDComment, server) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("вы не автор комментария")))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Удаляем комментарий
	comment, err := server.querier.DeleteComment(ctx, req.IDComment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

type getCommentRequest struct {
	IDComment int32 `uri:"id_comment" binding:"required,min=1"`
}

func (server *Server) getComment(ctx *gin.Context) {
	var req getCommentRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Возвращаем комментарий
	comment, err := server.querier.GetComment(ctx, req.IDComment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

type editCommentRequest struct {
	IDComment int32  `json:"id_comment" binding:"required,min=1"`
	Text      string `json:"text" binding:"required"`
}

func (server *Server) editComment(ctx *gin.Context) {
	var req editCommentRequest

	// Выявляем авторcтво
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if !isAuthorComment(payload.IDUser, req.IDComment, server) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("вы не автор комментария")))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Изменяем комментарий
	args := db.EditCommentParams{
		Text:      req.Text,
		IDComment: req.IDComment,
	}
	comment, err := server.querier.EditComment(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
}
