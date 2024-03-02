package api

import (
	"context"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createArticleRequest struct {
	Title   string  `json:"title"`
	Text    string  `json:"text"`
	Authors []int32 `json:"authors"`
}

func (proc *Process) createArticle(ctx *gin.Context) {
	var req createArticleRequest

	// Проверяем теги
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Создаём статью
	arg := db.CreateArticleParams{
		Title:   req.Title,
		Text:    req.Text,
		Authors: req.Authors,
	}
	article, err := proc.queries.CreateArticle(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)
}

type deleteArticleRequest struct {
	IDArticle int32 `uri:"id_article" binding:"required,min=1"`
}

func (proc *Process) deleteArticle(ctx *gin.Context) {
	var req deleteArticleRequest

	// Проверяем теги
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Удаляем статью
	article, err := proc.queries.DeleteArticle(context.Background(), req.IDArticle)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)
}

type getArticleRequest struct {
	IDArticle int32 `uri:"id_article" binding:"required,min=1"`
}

func (proc *Process) getArticle(ctx *gin.Context) {
	var req getArticleRequest

	// Проверяем теги
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем статью
	article, err := proc.queries.GetArticle(context.Background(), req.IDArticle)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// Да громоздко.
type getManyArticlesRequest struct {
	IDArticle  bool  `form:"id_article"`
	Evaluation bool  `form:"evaluation"`
	Comments   bool  `form:"comments"`
	Authors    bool  `form:"authors"`
	Title      bool  `form:"title"`
	Text       bool  `form:"text"`
	EditedAt   bool  `form:"edited_at"`
	CreatedAt  bool  `form:"created_at"`
	PageNum    int32 `form:"Offset" binding:"required,min=1"`
	PageSize   int32 `form:"Limit" binding:"required,min=1"`
}

func (proc *Process) getManySortedArticles(ctx *gin.Context) {
	var req getManyArticlesRequest

	// Проверяем теги
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем статьи
	arg := db.GetManySortedArticlesParams{
		IDArticle:  req.IDArticle,
		Evaluation: req.Evaluation,
		Comments:   req.Comments,
		Authors:    req.Authors,
		Title:      req.Title,
		Text:       req.Text,
		EditedAt:   req.EditedAt,
		CreatedAt:  req.CreatedAt,
		Limit:      req.PageSize,
		Offset:     (req.PageNum - 1) * req.PageSize,
	}
	articles, err := proc.queries.GetManySortedArticles(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

// Про прошлое "громоздко" я пошутил...
type getManySortedArticlesWithAttributesRequest struct {
	SelectByEditedAt   pgtype.Timestamp `form:"select_by_edited_at"`
	SelectByTitle      string           `form:"select_by_title"`
	SelectByText       string           `form:"select_by_text"`
	SelectByComments   []int32          `form:"select_by_comments"`
	SelectByAuthors    []int32          `form:"select_by_authors"`
	SelectByEvaluation int32            `form:"select_by_evaluation"`
	SortedByIDArticle  bool             `form:"sorted_by_id_article"`
	SortedByEvaluation bool             `form:"sorted_by_evaluation"`
	SortedByComments   bool             `form:"sorted_by_comments"`
	SortedByAuthors    bool             `form:"sorted_by_authors"`
	SortedByTitle      bool             `form:"sorted_by_title"`
	SortedByText       bool             `form:"sorted_by_text"`
	SortedByEditedAt   bool             `form:"sorted_by_edited_at"`
	SortedByCreatedAt  bool             `form:"sorted_by_created_at"`
	Offset             int32            `form:"Offset"`
	Limit              int32            `form:"Limit"`
}

func (proc *Process) getManySortedArticlesWithAttributes(ctx *gin.Context) {
	var req getManySortedArticlesWithAttributesRequest

	// Проверяем теги
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем статьи
	arg := db.GetManySortedArticlesWithAttributeParams{
		SelectByEditedAt:   req.SelectByEditedAt,
		SelectByTitle:      req.SelectByTitle,
		SelectByText:       req.SelectByText,
		SelectByComments:   req.SelectByComments,
		SelectByAuthors:    req.SelectByAuthors,
		SelectByEvaluation: req.SelectByEvaluation,
		SortedByIDArticle:  req.SortedByIDArticle,
		SortedByEvaluation: req.SortedByEvaluation,
		SortedByComments:   req.SortedByComments,
		SortedByAuthors:    req.SortedByAuthors,
		SortedByTitle:      req.SortedByTitle,
		SortedByText:       req.SortedByText,
		SortedByEditedAt:   req.SortedByEditedAt,
		SortedByCreatedAt:  req.SortedByCreatedAt,
		Offset:             req.Offset,
		Limit:              req.Limit,
	}
	articles, err := proc.queries.GetManySortedArticlesWithAttribute(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

// Надо разделить данные которые получаем с url и данные которые получаем с uri
type idEditArticleRequest struct {
	IDArticle int32 `form:"id_article" binding:"required"`
}
type editArticleRequest struct {
	Title      string           `form:"title"`
	Text       string           `form:"text"`
	Comments   []int32          `form:"comments"`
	Authors    []int32          `form:"authors"`
	Evaluation int32            `form:"evaluation"`
	EditedAt   pgtype.Timestamp `form:"edited_at"`
}

func (proc *Process) editArticle(ctx *gin.Context) {
	var artID idEditArticleRequest
	var req editArticleRequest

	// Проверяем теги
	if err := ctx.ShouldBindUri(&artID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Изменяем статью
	arg := db.EditArticleParamParams{
		IDArticle:  artID.IDArticle,
		Title:      req.Title,
		Text:       req.Text,
		Comments:   req.Comments,
		Authors:    req.Authors,
		Evaluation: req.Evaluation,
		EditedAt:   req.EditedAt,
	}

	editedArticle, err := proc.queries.EditArticleParam(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, editedArticle)
}
