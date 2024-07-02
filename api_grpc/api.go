package api_grpc

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/pb"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
)

// Server Обрабатываем запросы от API по gRPS
type Server struct {
	pb.UnimplementedEducationServer
	queries    *db.Queries
	tokenMaker token.Maker
	config     tools.Config
}

// NewServer Новый gRPC процесс для обработки запросов
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

	return server, nil
}
