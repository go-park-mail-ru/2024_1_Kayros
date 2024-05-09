package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"2024_1_kayros/microservices/user/internal/repo"
	"2024_1_kayros/microservices/user/internal/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/config"
	"2024_1_kayros/internal/repository/minios3"
	"2024_1_kayros/services/minio"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.UserGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice user doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice user responds on port %d", cfg.UserGrpcServer.Port))

	// init grpc server
	server := grpc.NewServer()
	// init services for server work
	postgreDB := postgres.Init(cfg, logger)
	minioClient := minio.Init(cfg, logger)

	repoUser := repo.NewLayer(postgreDB)
	repoMinio := minios3.NewRepoLayer(minioClient)
	user.RegisterUserManagerServer(server, usecase.NewLayer(repoUser, repoMinio))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", cfg.UserGrpcServer.Host, cfg.UserGrpcServer.Port), zap.String("error", err.Error()))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice user has shut down")
	os.Exit(0)
}
