package api

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
	"github.com/ZenSam7/Education/worker"
)

// Server Обрабатываем запросы от API по gRPS
type Server struct {
	pb.UnimplementedEducationServer
	queries         *db.Queries
	tokenMaker      token.Maker
	config          tools.Config
	taskDistributor worker.TaskDistributor
}

// NewServer Новый gRPC процесс для обработки запросов
func NewServer(config tools.Config, queries *db.Queries, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		queries:         queries,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
