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
	querier         db.Querier
	tokenMaker      token.Maker
	config          tools.Config
	taskDistributor worker.TaskDistributor
}

// NewServer Новый gRPC процесс для обработки запросов (используем Paseto)
func NewServer(config tools.Config, querier db.Querier, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		querier:         querier,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
