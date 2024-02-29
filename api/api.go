package api

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

// NewProcess Новый HTTP процесс для обработки и роутер
func NewProcess(queries *db.Queries) *Process {
	proc := &Process{queries: queries}
	router := gin.Default()

	router.POST("/user", proc.createUser)
	// ":id_user" Даём gin понять что нам нужен парамерт URI id_user
	router.GET("/user/:id_user", proc.getUser)
	router.GET("/user", proc.getManyUsers)

	proc.router = router
	return proc
}

// errorResponse Преврящаем ошибку в нужный объект чтобы использовать его в gin
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
