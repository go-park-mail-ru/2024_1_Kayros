package main

import (
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/microservices/restaurants/internal/repo"
	"2024_1_kayros/microservices/restaurants/internal/usecase"
	rest "2024_1_kayros/microservices/restaurants/proto"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.RestGrpcServer.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("The server cannot be started.\n%v", err)
	} else {
		log.Printf("The server listen port %d", cfg.RestGrpcServer.Port)
	}

	server := grpc.NewServer()

	postgreDB := postgres.Init(cfg, logger)
	repoRest := repo.NewRestLayer(postgreDB)
	rest.RegisterRestWorkerServer(server, usecase.NewRestLayer(repoRest))
	err = server.Serve(lis)
	if err != nil {
		log.Printf("error in serving server on port %d -  %s", cfg.RestGrpcServer.Port, err)
	}
}
