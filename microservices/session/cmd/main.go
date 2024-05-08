package main

import (
	"fmt"
	"log"
	"net"

	"2024_1_kayros/microservices/session/internal/repo"
	"2024_1_kayros/microservices/session/internal/usecase"
	sessionv1 "2024_1_kayros/microservices/session/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/redis"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.SessionGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		errMsg := fmt.Sprintf("The server cannot be started.\n%v", err)
		logger.Fatal(errMsg)
	}
	infoMsg := fmt.Sprintf("The session server listens port %d", cfg.SessionGrpcServer.Port)
	logger.Info(infoMsg)

	server := grpc.NewServer()
	redisSession := redis.Init(cfg, logger, cfg.Redis.DatabaseSession)
	redisCsrf := redis.Init(cfg, logger, cfg.Redis.DatabaseCsrf)

	repoSession := repo.NewLayer(redisSession)
	repoCsrf := repo.NewLayer(redisCsrf)
	
	sessionv1.RegisterSessionManagerServer(server, usecase.NewLayer(repoCsrf, repoSession, &cfg.Redis))
	err = server.Serve(conn)
	if err != nil {
		log.Fatalf("error in serving server on port %d -  %s", cfg.SessionGrpcServer.Port, err)
	}
}
