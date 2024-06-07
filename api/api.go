package api

// Что должно делать наже API при каки-либо обращениях к нему

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/gin-gonic/gin"
)

// Process Обрабатываем запросы от API
type Process struct {
	queries    *db.Queries
	router     *gin.Engine
	tokenMaker token.Maker
	config     tools.Config
}

// Run Начинаем прослушивать запросы к API
func (proc *Process) Run(address string) error {
	return proc.router.Run(address)
}

// NewProcess Новый HTTP процесс для обработки запросов и роутер (который просто
// вызывает определёную функцию при каком-либо запросе на конкретный URI)
func NewProcess(config tools.Config, queries *db.Queries) (*Process, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	proc := &Process{
		queries:    queries,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := gin.Default()
	proc.setupRouter(router)

	return proc, nil
}

// setupRouter Устанавливаем все возможные url для обработки, а также делим
// их для авторизованных и не авторизованных пользователей
func (proc *Process) setupRouter(router *gin.Engine) {
	// Добавляем сайты только для авторизованных пользователей ("/" - общий префикс)
	authRouter := router.Group("/").Use(authMiddleware(proc.tokenMaker))

	// Обрабатываем запросы для действий с пользователями:
	router.POST("/user", proc.createUser)
	router.POST("/user/login", proc.loginUser)
	// ":id_user" Даём gin понять что нам нужен парамерт URI id_user
	router.GET("/user/:id_user", proc.getUser)
	router.GET("/user/list", proc.getManySortedUsers)
	authRouter.PATCH("/user/", proc.editUserParam)
	authRouter.DELETE("/user/", proc.deleteUser)

	// Обрабатываем запросы для действий со статьями:
	authRouter.POST("/article", proc.createArticle)
	authRouter.DELETE("/article/:id_article", proc.deleteArticle)
	router.GET("/article/:id_article", proc.getArticle)
	router.GET("/article/list", proc.getManySortedArticles)
	router.GET("/article/comments/:id_article", proc.getCommentsOfArticle)
	authRouter.PATCH("/article/:id_article", proc.editArticle)
	router.GET("/article/search", proc.getManySortedArticlesWithAttributes)

	// Обрабатываем запросы для действий с комментариями:
	authRouter.POST("/comment", proc.createComment)
	router.GET("/comment/:id_comment", proc.getComment)
	authRouter.PATCH("/comment/:id_comment", proc.editComment)
	authRouter.DELETE("/comment/:id_comment", proc.deleteComment)

	proc.router = router
}

// errorResponse Превpящаем ошибку в нужный объект чтобы использовать его в gin
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
