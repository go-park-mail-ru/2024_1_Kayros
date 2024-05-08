package main

import (
	"fmt"
	"log"
	"net"

	rest "2024_1_kayros/microservices/comment/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.AuthGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		errMsg := fmt.Sprintf("The server cannot be started.\n%v", err)
		logger.Fatal(errMsg)
	}
	infoMsg := fmt.Sprintf("The authentication server listens port %d", cfg.CommentGrpcServer.Port)
	logger.Info(infoMsg)

	server := grpc.NewServer()
	postgreDB := postgres.Init(cfg, logger)
	repoComment := repo.NewCommentLayer(postgreDB)
	rest.RegisterCommentWorkerServer(server, usecase.NewCommentLayer(repoComment))
	err = server.Serve(conn)
	if err != nil {
		log.Fatalf("error in serving server on port %d -  %s", cfg.CommentGrpcServer.Port, err)
	}
}
