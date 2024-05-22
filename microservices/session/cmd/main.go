package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/microservices/session/internal/repo"
	"2024_1_kayros/microservices/session/internal/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/services/redis"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.SessionGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice session doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice session responds on port %d", cfg.SessionGrpcServer.Port))

	server := grpc.NewServer()
	redisSession := redis.Init(cfg, logger, cfg.Redis.DatabaseSession)
	redisCsrf := redis.Init(cfg, logger, cfg.Redis.DatabaseCsrf)

	repoSession := repo.NewLayer(redisSession)
	repoCsrf := repo.NewLayer(redisCsrf)

	session.RegisterSessionManagerServer(server, usecase.NewLayer(repoCsrf, repoSession, logger, &cfg.Redis))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", cfg.SessionGrpcServer.Host, cfg.SessionGrpcServer.Port), zap.String("error", err.Error()))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice session has shut down")
	os.Exit(0)
}