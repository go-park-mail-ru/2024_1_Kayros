package main

import (
	"log"
	"net"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/microservices/restaurants/internal/repo"
	"2024_1_kayros/microservices/restaurants/internal/usecase"
	rest "2024_1_kayros/microservices/restaurants/proto"
	"2024_1_kayros/services/postgres"
)

const PORT = 8081

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		log.Println("The server cannot be started.\n%v", err)
	} else {
		log.Println("The server listen port", PORT)
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
		log.Printf("error in serving server on port %d -  %s", PORT, err)
	}
}
