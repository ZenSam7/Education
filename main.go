package main

import (
	"context"
	"github.com/ZenSam7/Education/api"
	"github.com/ZenSam7/Education/api_grpc"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/tools"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {
	config := tools.LoadConfig(".")
	queries, closeConn := db.GetQueries()
	defer closeConn()

	go runGatewayServer(config, queries)
	runGrpcServer(config, queries)
}

// runGatewayServer Сервер на gRPC, но с поддержкой HTTP
func runGatewayServer(config tools.Config, queries *db.Queries) {
	server, err := api_grpc.NewServer(config, queries)
	if err != nil {
		log.Fatal("Ошибка в создании сервера:", err.Error())
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
		log.Fatal("не получилось поднять gRPC Gateway сервер:", err)
	}

	// Прослушиваем все адреса (т.к. корневой узел)
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("не получилось создать listener:", err.Error())
	}

	log.Printf("gRPC сервер поднят на %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("не получилось поднять gRPC Gateway сервер:", err)
	}
}

// runGrpcServer Стандартный сервер на gRPC
func runGrpcServer(config tools.Config, queries *db.Queries) {
	server, err := api_grpc.NewServer(config, queries)
	if err != nil {
		log.Fatal("Ошибка в создании сервера:", err.Error())
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterEducationServer(grpcServer, server)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("не получилось создать listener:", err.Error())
	}

	log.Printf("gRPC сервер поднят на %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("не получилось поднять gRPC сервер:", err)
	}
}

// runGrpcServer Стандартный сервер на Gin
func runGinServer(config tools.Config, queries *db.Queries) {
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatal("Ошибка в создании сервера:", err.Error())
	}

	if err := server.Run(config.HttpServerAddress); err != nil {
		log.Fatal("Не получилось поднять Gin сервер:", err)
	}
}
