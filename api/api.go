package api

import (
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/redis/cache"
	"github.com/ZenSam7/Education/redis/worker"
	"github.com/ZenSam7/Education/token"
	"github.com/ZenSam7/Education/tools"
)

// Server Обрабатываем запросы от API по gRPS
type Server struct {
	pb.UnimplementedEducationServer
	querier         db.Querier
	replicaConn     db.Querier
	tokenMaker      token.Maker
	config          tools.Config
	taskDistributor worker.TaskDistributor
	cacher          cache.Cacher
}

// NewServer Новый gRPC процесс для обработки запросов (используем Paseto)
func NewServer(
	querier db.Querier,
	replicaConn db.Querier,
	config tools.Config,
	tokenMaker token.Maker,
	taskDistributor worker.TaskDistributor,
	cacher cache.Cacher,
) *Server {

	server := &Server{
		querier:         querier,
		replicaConn:     replicaConn,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
		cacher:          cacher,
	}

	return server
}
