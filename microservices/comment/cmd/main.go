package main

import (
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/microservices/comment/internal/repo"
	"2024_1_kayros/microservices/comment/internal/usecase"
	rest "2024_1_kayros/microservices/comment/proto"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.CommentGrpcServer.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("The server cannot be started.\n%v", err)
	} else {
		log.Printf("The server listen port %d", cfg.CommentGrpcServer.Port)
	}

	server := grpc.NewServer()

	postgreDB := postgres.Init(cfg, logger)
	repoComment := repo.NewCommentLayer(postgreDB)
	rest.RegisterCommentWorkerServer(server, usecase.NewCommentLayer(repoComment))
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("error in serving server on port %d -  %s", cfg.CommentGrpcServer.Port, err)
	}
}
