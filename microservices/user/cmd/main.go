package main

import (
	"fmt"
	"log"
	"net"

	"2024_1_kayros/microservices/user/internal/repo"
	"2024_1_kayros/microservices/user/internal/usecase"
	"2024_1_kayros/microservices/user/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/repository/minios3"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/minio"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.UserGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		errMsg := fmt.Sprintf("The server cannot be started.\n%v", err)
		logger.Fatal(errMsg)
	}
	infoMsg := fmt.Sprintf("The user server listens port %d", cfg.UserGrpcServer.Port)
	logger.Info(infoMsg)

	server := grpc.NewServer()
	postgreDB := postgres.Init(cfg, logger)
	minioClient := minio.Init(cfg, logger)

	repoUser := repo.NewLayer(postgreDB)
	repoMinio := minios3.NewRepoLayer(minioClient)
	userv1.RegisterUserManagerServer(server, usecase.NewLayer(repoUser, repoMinio))
	err = server.Serve(conn)
	if err != nil {
		log.Fatalf("error in serving server on port %d -  %s", cfg.UserGrpcServer.Port, err)
	}
}
