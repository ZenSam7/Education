package api

// Что должно делать наже API при каки-либо обращениях к нему

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Process Обрабатываем запросы от API
type Process struct {
	queries *db.Queries
	router  *gin.Engine
}

// Run Начинаем прослушивать запросы к API
func (proc *Process) Run(address string) error {
	return proc.router.Run(address)
}

// NewProcess Новый HTTP процесс для обработки запросов и роутер (который просто
// вызывает определёную функцию при каком-либо запросе на конкретный URI)
func NewProcess(queries *db.Queries) *Process {
	proc := &Process{queries: queries}
	router := gin.Default()

	// Как обрабатываем запросы для действий с пользователями:
	router.POST("/user", proc.createUser)
	// ":id_user" Даём gin понять что нам нужен парамерт URI id_user
	router.GET("/user/:id_user", proc.getUser)
	router.GET("/user", proc.getManyUsers)
	router.PATCH("/user/:id_user", proc.editUserParam)
	router.DELETE("/user/:id_user", proc.deleteUser)

	// Как обрабатываем запросы для действий со статьями:
	router.POST("/article", proc.createArticle)
	router.DELETE("/article/:id_article", proc.deleteArticle)
	router.GET("/article/:id_article", proc.getArticle)
	router.GET("/article", proc.getManySortedArticles)
	router.PATCH("/article/:id_article", proc.editArticle)
	router.GET("/article/search", proc.getManySortedArticlesWithAttributes)

	// Как обрабатываем запросы для действий с комментариями:
	router.POST("/comment", proc.createComment)
	router.GET("/comment/:id_comment", proc.getComment)
	router.GET("/comment", proc.getCommentsOfArticle)
	router.PATCH("/comment/:id_comment", proc.editComment)
	router.DELETE("/comment/:id_comment", proc.deleteComment)

	proc.router = router
	return proc
}

// errorResponse Преврящаем ошибку в нужный объект чтобы использовать его в gin
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
