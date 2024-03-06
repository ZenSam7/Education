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
	IDArticle  bool  `json:"id_article"`
	Evaluation bool  `json:"evaluation"`
	Comments   bool  `json:"comments"`
	Authors    bool  `json:"authors"`
	Title      bool  `json:"title"`
	Text       bool  `json:"text"`
	EditedAt   bool  `json:"edited_at"`
	CreatedAt  bool  `json:"created_at"`
	PageNum    int32 `json:"page_num" binding:"required,min=1"`
	PageSize   int32 `json:"page_size" binding:"required,min=1"`
}

func (proc *Process) getManySortedArticles(ctx *gin.Context) {
	var req getManyArticlesRequest

	// Проверяем теги
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
	SelectEditedAt   pgtype.Timestamp `json:"select_edited_at"`
	SelectTitle      string           `json:"select_title"`
	SelectText       string           `json:"select_text"`
	SelectComments   []int32          `json:"select_comments"`
	SelectAuthors    []int32          `json:"select_authors"`
	SelectEvaluation int32            `json:"select_evaluation"`
	SortedIDArticle  bool             `json:"sorted_id_article"`
	SortedEvaluation bool             `json:"sorted_evaluation"`
	SortedComments   bool             `json:"sorted_comments"`
	SortedAuthors    bool             `json:"sorted_authors"`
	SortedTitle      bool             `json:"sorted_title"`
	SortedText       bool             `json:"sorted_text"`
	SortedEditedAt   bool             `json:"sorted_edited_at"`
	SortedCreatedAt  bool             `json:"sorted_created_at"`
	PageNum          int32            `json:"page_num" binding:"required,min=1"`
	PageSize         int32            `json:"page_size" binding:"required,min=1"`
}

func (proc *Process) getManySortedArticlesWithAttributes(ctx *gin.Context) {
	var req getManySortedArticlesWithAttributesRequest

	// Проверяем теги
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Получаем статьи
	arg := db.GetManySortedArticlesWithAttributeParams{
		SelectTitle:      req.SelectTitle,
		SelectText:       req.SelectText,
		SelectEvaluation: req.SelectEvaluation,
		SortedIDArticle:  req.SortedIDArticle,
		SortedEvaluation: req.SortedEvaluation,
		SortedComments:   req.SortedComments,
		SortedAuthors:    req.SortedAuthors,
		SortedTitle:      req.SortedTitle,
		SortedText:       req.SortedText,
		SortedEditedAt:   req.SortedEditedAt,
		SortedCreatedAt:  req.SortedCreatedAt,
		Limit:            req.PageSize,
		Offset:           (req.PageNum - 1) * req.PageSize,
	}
	articles, err := proc.queries.GetManySortedArticlesWithAttribute(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

// Надо разделить данные которые получаем с url и данные которые получаем с uri
type editArticleRequest struct {
	IDArticle  int32            `json:"id_article" binding:"required,min=1"`
	Title      string           `json:"title"`
	Text       string           `json:"text"`
	Comments   []int32          `json:"comments"`
	Authors    []int32          `json:"authors"`
	Evaluation int32            `json:"evaluation"`
	EditedAt   pgtype.Timestamp `json:"edited_at"`
}

func (proc *Process) editArticle(ctx *gin.Context) {
	var req editArticleRequest

	// Проверяем теги
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Изменяем статью
	arg := db.EditArticleParams{
		IDArticle:  req.IDArticle,
		Title:      req.Title,
		Text:       req.Text,
		Comments:   req.Comments,
		Authors:    req.Authors,
		Evaluation: req.Evaluation,
	}

	editedArticle, err := proc.queries.EditArticle(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, editedArticle)
}
