package main

import (
	"fmt"
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
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("can not listen port 8081: %s", err)
	}

	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)
	cfg := config.NewConfig(logger)

	server := grpc.NewServer()

	postgreDB := postgres.Init(cfg, logger)
	repoRest := repo.NewRestLayer(postgreDB)
	rest.RegisterRestWorkerServer(server, usecase.NewRestLayer(repoRest))
	err = server.Serve(lis)
	if err != nil {
		fmt.Println("error in serving server on port 8081 %s", err)
	}
}
