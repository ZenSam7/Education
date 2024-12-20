package api_gin

import (
	"context"
	"errors"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type createArticleRequest struct {
	Title   string  `json:"title"`
	Text    string  `json:"text"`
	Authors []int32 `json:"authors"`
}

func (server *Server) createArticle(ctx *gin.Context) {
	var req createArticleRequest

	// Делаем операцию только для авторизованного пользователя
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)

	// Проверяем теги
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Создаём статью
	arg := db.CreateArticleParams{
		Title:   req.Title,
		Text:    req.Text,
		Authors: []int32{payload.IDUser},
	}
	article, err := server.querier.CreateArticle(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// userIsAuthorArticle Проверяем является ли пользователь автором статьи
func isAuthorArticle(IDArticle, IDUser int32, server *Server) bool {
	targetArticle, _ := server.querier.GetArticle(context.Background(), IDArticle)
	for _, authorID := range targetArticle.Authors {
		if authorID == IDUser {
			return true
		}
	}
	return false
}

type deleteArticleRequest struct {
	IDArticle int32 `uri:"id_article" binding:"required,min=1"`
}

func (server *Server) deleteArticle(ctx *gin.Context) {
	var req deleteArticleRequest

	// Делаем операцию только для создателя статьи
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)

	if !isAuthorArticle(req.IDArticle, payload.IDUser, server) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("вы не являетесь автором статьи")))
		return
	}

	// Проверяем теги
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Удаляем статью
	article, err := server.querier.DeleteArticle(ctx, req.IDArticle)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)
}

type getArticleRequest struct {
	IDArticle int32 `uri:"id_article" binding:"required,min=1"`
}

func (server *Server) getArticle(ctx *gin.Context) {
	var req getArticleRequest

	// Проверяем теги
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Попытка получения данных из кеша
	cacheKey := fmt.Sprintf("article:%d", req.IDArticle)
	var cachedResponse db.Article
	if err := server.cacher.GetCache(ctx, cacheKey, &cachedResponse); err == nil {
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	// Получаем статью
	article, err := server.querier.GetArticle(ctx, req.IDArticle)
	if err != nil {
		// Если не получилось с главной бд, пытаемся с репликой
		article, err = server.replicaConn.GetArticle(ctx, req.IDArticle)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)

	if err := server.cacher.SetCache(ctx, cacheKey, article); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("не удалось сохранить данные в кеш: %s", err)))
		return
	}
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

func (server *Server) getManySortedArticles(ctx *gin.Context) {
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
	articles, err := server.querier.GetManySortedArticles(ctx, arg)
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

func (server *Server) getManySortedArticlesWithAttributes(ctx *gin.Context) {
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
	articles, err := server.querier.GetManySortedArticlesWithAttribute(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

// Надо разделить данные которые получаем с url и данные которые получаем с uri
type editArticleRequest struct {
	IDArticle int32            `json:"id_article" binding:"required,min=1"`
	Title     string           `json:"title"`
	Text      string           `json:"text"`
	Comments  []int32          `json:"comments"`
	Authors   []int32          `json:"authors"`
	EditedAt  pgtype.Timestamp `json:"edited_at"`
}

func (server *Server) editArticle(ctx *gin.Context) {
	var req editArticleRequest

	// Делаем операцию только для создателя статьи
	payload := ctx.MustGet(authPayloadKey).(*token.Payload)

	if !isAuthorArticle(req.IDArticle, payload.IDUser, server) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("вы не являетесь автором статьи")))
		return
	}

	// Проверяем теги
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Изменяем статью
	arg := db.EditArticleParams{
		IDArticle: req.IDArticle,
		Title:     req.Title,
		Text:      req.Text,
		Comments:  req.Comments,
		Authors:   req.Authors,
	}

	editedArticle, err := server.querier.EditArticle(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, editedArticle)
}
