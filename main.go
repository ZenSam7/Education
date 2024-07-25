package main

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/api"
	"github.com/ZenSam7/Education/api_grpc"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/tools"
	"github.com/ZenSam7/Education/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
)

func main() {
	config := tools.LoadConfig(".")
	queries, closeConn := db.GetQueries()
	defer closeConn() // (На самом деле оно не вызывается)
	tools.MakeLogger()

	runDBMigration(config)
	redisOpt, taskDistributor := makeRedis(config)

	go startTaskProcessor(redisOpt, queries)
	go runGatewayServer(config, queries, taskDistributor)
	runGrpcServer(config, queries, taskDistributor)
}

// startTaskProcessor Запускаем обработчик процессов
func startTaskProcessor(options asynq.RedisClientOpt, queries *db.Queries) {
	processor := worker.NewRedisTaskProcessor(options, queries)
	err := processor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("task-процессор не хочет создаваться")
	}
}

// makeRedis Запускаем сервер редиски
func makeRedis(config tools.Config) (asynq.RedisClientOpt, worker.TaskDistributor) {
	options := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	return options, worker.NewRedisTaskDistributor(options)
}

// runDBMigration Запускаем миграции через Go
func runDBMigration(config tools.Config) {
	migration, err := migrate.New(config.MigrationUrl, fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/education?sslmode=%s",
		config.DBUserName,
		config.DBPassword,
		config.DBHost,
		config.DBSSLMode,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось создать миграцию")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("не получилось поднять миграцию")
	}

	log.Info().Msg("миграция завершена")
}

// runGatewayServer Сервер на gRPC, но с поддержкой HTTP
func runGatewayServer(config tools.Config, queries *db.Queries, taskDistributor worker.TaskDistributor) {
	server, err := api_grpc.NewServer(config, queries, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка в создании сервера")
	}

	// Важная штука, которая не изменяет названия в json'е (названия остаются в snake_style)
	jsobOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsobOptions)

	err = pb.RegisterEducationHandlerServer(context.Background(), grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось поднять gRPC Gateway сервер")
	}

	// Прослушиваем все адреса (т.к. корневой узел)
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось создать listener")
	}
	log.Info().Msgf("gRPC Gateway сервер поднят на %s", listener.Addr().String())

	// Создаём специальный логгер для http
	handler := tools.HttpLogger(mux)

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось поднять gRPC Gateway сервер")
	}
}

// runGrpcServer Стандартный сервер на gRPC
func runGrpcServer(config tools.Config, queries *db.Queries, taskDistributor worker.TaskDistributor) {
	server, err := api_grpc.NewServer(config, queries, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка в создании сервера")
	}

	// Настраиваем логгер
	lggr := grpc.UnaryInterceptor(tools.GrpcLogger)

	grpcServer := grpc.NewServer(lggr)
	reflection.Register(grpcServer)
	pb.RegisterEducationServer(grpcServer, server)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось создать listener")
	}

	log.Info().Msgf("gRPC сервер поднят на %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось поднять gRPC сервер")
	}
}

// runGrpcServer Стандартный сервер на Gin
func runGinServer(config tools.Config, queries *db.Queries) {
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка в создании сервера")
	}

	if err := server.Run(config.HttpServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Не получилось поднять Gin сервер:")
	}
}
