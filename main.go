package main

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/api"
	"github.com/ZenSam7/Education/api_gin"
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
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var config tools.Config
var queries *db.Queries
var closeConn func()

// interruptSignals список ошибок, которые не дадут серверу упасть мнговенно; когда возникнет одна из этих ошибок,
// сервер завершит обработку всех текущих запросов, а уже потом ляжет
var interruptSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
	syscall.SIGTERM,
	syscall.SIGINT,
}

var ctx context.Context
var waitErr *errgroup.Group

func main() {
	// Настройки
	config = tools.LoadConfig()
	tools.MakeLogger()

	// Бд
	queries, closeConn = db.MakeQueries()
	defer closeConn()
	runDBMigration()
	redisOpt, taskDistributor := makeRedis()

	// Захватываем ошибки во время работы серверов
	notifyCtx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	waitErr, ctx = errgroup.WithContext(notifyCtx)

	// Запускаем сервера
	startTaskProcessor(redisOpt)
	runHttpGatewayServer(taskDistributor)
	runGrpcServer(taskDistributor)
	//runGinServer()

	if err := waitErr.Wait(); err != nil {
		log.Fatal().Err(err).Msg("сервер лёг")
		stop()
	}
}

// startTaskProcessor Запускаем обработчик процессов
func startTaskProcessor(options asynq.RedisClientOpt) {
	processor := worker.NewRedisTaskProcessor(options, queries)

	// Запускаем сервер конкурентно, и, если что, захватываем ошибку в waitErr
	waitErr.Go(func() error {
		err := processor.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("task-процессор не хочет создаваться ((")
		}
		return err
	})
	// Мягонько и аккуратненько останавливаем сервер в случае ошибочки
	waitErr.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("останавливаем TaskProcessor")
		processor.Shutdown()
		return nil
	})
}

// makeRedis Запускаем сервер редиски
func makeRedis() (asynq.RedisClientOpt, worker.TaskDistributor) {
	options := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	return options, worker.NewRedisTaskDistributor(options)
}

// runDBMigration Запускаем миграции через Go
func runDBMigration() {
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

// runHttpGatewayServer Сервер на gRPC, но с поддержкой HTTP
func runHttpGatewayServer(taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, queries, taskDistributor)
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
		log.Fatal().Err(err).Msg("не получилось поднять HTTP Gateway сервер")
	}

	// Прослушиваем все адреса (т.к. корневой узел)
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Создаём специальный логгер для http
	handler := tools.HttpLogger(mux)

	httpGatewayServer := &http.Server{
		Addr:    config.HttpServerAddress,
		Handler: handler,
	}

	// Запускаем сервер конкурентно, и, если что, захватываем ошибку в waitErr
	waitErr.Go(func() error {
		log.Info().Msgf("HTTP Gateway сервер поднят на %s", httpGatewayServer.Addr)

		if err = httpGatewayServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return nil
			}

			log.Fatal().Err(err).Msg("не получилось поднять HTTP Gateway сервер")
		}
		return err
	})
	// Мягонько и аккуратненько останавливаем сервер в случае ошибочки
	waitErr.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("останавливаем HTTP Gateway")
		return httpGatewayServer.Shutdown(context.Background())
	})
}

// runGrpcServer Стандартный сервер на gRPC
func runGrpcServer(taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, queries, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка в создании сервера")
	}

	// Настраиваем логгер
	lggr := grpc.UnaryInterceptor(tools.GrpcLogger)

	// Регистрируем
	grpcServer := grpc.NewServer(lggr)
	reflection.Register(grpcServer)
	pb.RegisterEducationServer(grpcServer, server)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось создать listener")
	}

	log.Info().Msgf("gRPC сервер поднят на %s", listener.Addr().String())

	// Запускаем сервер конкурентно, и, если что, захватываем ошибку в waitErr
	waitErr.Go(func() error {
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("не получилось поднять gRPC сервер")
		}
		return err
	})
	// Мягонько и аккуратненько останавливаем сервер в случае ошибочки
	waitErr.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("останавливаем gRPC")
		grpcServer.GracefulStop()
		return nil
	})
}

// runGrpcServer Стандартный сервер на Gin
func runGinServer() {
	server, err := api_gin.NewServer(config, queries)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка в создании сервера")
	}

	// Запускаем сервер конкурентно, и, если что, захватываем ошибку в waitErr
	waitErr.Go(func() error {
		if err := server.Run(config.HttpServerAddress); err != nil {
			log.Fatal().Err(err).Msg("Не получилось поднять Gin сервер:")
		}
		return err
	})
	// Убиваем сервер нахуй
	return
}
