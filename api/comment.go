package api

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getCommentsOfArticleRequest struct {
	IDArticle int32 `json:"id_article" binding:"required,min=1"`
	PageNum   int32 `json:"page_num" binding:"required,min=1"`
	PageSize  int32 `json:"page_size" binding:"required,min=1"`
}

func (proc *Process) getCommentsOfArticle(c *gin.Context) {
	var req getCommentsOfArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Возвращаем комментарии
	args := db.GetCommentsOfArticleParams{
		IDArticle: req.IDArticle,
		Limit:     req.PageSize,
		Offset:    (req.PageNum - 1) * req.PageSize,
	}
	comments, err := proc.queries.GetCommentsOfArticle(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comments)
}

type createCommentRequest struct {
	IDArticle int32  `json:"id_article" binding:"required,min=1"`
	Text      string `json:"text" binding:"required"`
	FromUser  int32  `json:"from_user" binding:"required"`
}

func (proc *Process) createComment(c *gin.Context) {
	var req createCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Создаем комментарий
	args := db.CreateCommentParams{
		IDArticle: req.IDArticle,
		Text:      req.Text,
		FromUser:  req.FromUser,
	}
	comment, err := proc.queries.CreateComment(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}

type deleteCommentRequest struct {
	IDComment int32 `uri:"id_comment" binding:"required,min=1"`
}

func (proc *Process) deleteComment(c *gin.Context) {
	var req deleteCommentRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Удаляем комментарий
	comment, err := proc.queries.DeleteComment(c, req.IDComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}

type getCommentRequest struct {
	IDComment int32 `uri:"id_comment" binding:"required,min=1"`
}

func (proc *Process) getComment(c *gin.Context) {
	var req getCommentRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Возвращаем комментарий
	comment, err := proc.queries.GetComment(c, req.IDComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}

type editCommentRequest struct {
	IDComment int32  `json:"id_comment" binding:"required,min=1"`
	Text      string `json:"text" binding:"required"`
}

func (proc *Process) editComment(c *gin.Context) {
	var req editCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Изменяем комментарий
	args := db.EditCommentParams{
		Text:      req.Text,
		IDComment: req.IDComment,
	}
	comment, err := proc.queries.EditComment(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}
