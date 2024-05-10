package main

import (
	"2024_1_kayros/gen/go/food"
	"2024_1_kayros/microservices/food/internal/repo"
	"2024_1_kayros/microservices/food/internal/usecase"
	"2024_1_kayros/services/postgres"
	"fmt"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.RestGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice restaurant doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice restaurant responds on port %d", cfg.RestGrpcServer.Port))

	// init grpc server 
	server := grpc.NewServer()
	// init services for server work
	postgreDB := postgres.Init(cfg, logger)
	// register contract 
	repoUser := repo.NewLayer(postgreDB)
	food.RegisterFoodManagerServer(server, usecase.NewLayer(repoUser))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", cfg.RestGrpcServer.Host, cfg.RestGrpcServer.Port), zap.String("error", err.Error()))
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice authorization has shut down")
	os.Exit(0)
}
