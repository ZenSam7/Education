package api_gin

// Что должно делать наже API при каки-либо обращениях к нему

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/gin-gonic/gin"
)

// Server Обрабатываем запросы от API
type Server struct {
	queries    *db.Queries
	router     *gin.Engine
	tokenMaker token.Maker
	config     tools.Config
}

// Run Начинаем прослушивать запросы к API по HTTP
func (server *Server) Run(address string) error {
	return server.router.Run(address)
}

// NewServer Новый HTTP процесс для обработки запросов и роутер (который просто
// вызывает определёную функцию при каком-либо запросе на конкретный URI)
func NewServer(config tools.Config, queries *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		queries:    queries,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := gin.Default()
	server.setupRouter(router)

	return server, nil
}

// setupRouter Устанавливаем все возможные url для обработки, а также делим
// их для авторизованных и не авторизованных пользователей
func (server *Server) setupRouter(router *gin.Engine) {
	// Добавляем сайты только для авторизованных пользователей ("/" - общий префикс)
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Обрабатываем запросы для действий с пользователями:
	router.POST("/user/register", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)
	// ":id_user" Даём gin понять что нам нужен парамерт URI id_user
	router.GET("/user/:id_user", server.getUser)
	router.GET("/user/list", server.getManySortedUsers)
	authRouter.PATCH("/user/", server.editUser)
	authRouter.DELETE("/user/", server.deleteUser)

	// Обрабатываем запросы для действий со статьями:
	authRouter.POST("/article", server.createArticle)
	authRouter.DELETE("/article/:id_article", server.deleteArticle)
	router.GET("/article/:id_article", server.getArticle)
	router.GET("/article/list", server.getManySortedArticles)
	router.GET("/article/comments/:id_article", server.getCommentsOfArticle)
	authRouter.PATCH("/article/:id_article", server.editArticle)
	router.GET("/article/search", server.getManySortedArticlesWithAttributes)

	// Обрабатываем запросы для действий с комментариями:
	authRouter.POST("/comment", server.createComment)
	router.GET("/comment/:id_comment", server.getComment)
	authRouter.PATCH("/comment/:id_comment", server.editComment)
	authRouter.DELETE("/comment/:id_comment", server.deleteComment)

	server.router = router
}

// errorResponse Превpящаем ошибку в нужный объект чтобы использовать его в gin
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
