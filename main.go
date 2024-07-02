package main

import (
	"github.com/ZenSam7/Education/api"
	"github.com/ZenSam7/Education/api_grpc"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/pb"
	"github.com/ZenSam7/Education/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	config := tools.LoadConfig(".")
	queries, closeConn := db.GetQueries()
	defer closeConn()

	runGrpcServer(config, queries)
}

func runGrpcServer(config tools.Config, queries *db.Queries) {
	server, err := api_grpc.NewServer(config, queries)
	if err != nil {
		log.Fatal("Ошибка в создании роутера:", err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEducationServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("не получилось создать listener:", err.Error())
	}

	log.Printf("gRPC сервер поднят на %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("не получилось поднять gRPC сервер:", err.Error())
	}
}

func runGinServer(config tools.Config, queries *db.Queries) {
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatal("Ошибка в создании роутера:", err.Error())
	}

	if err := server.Run(config.HttpServerAddress); err != nil {
		log.Fatal("Не получилось поднять Gin сервер:", err)
	}
}
