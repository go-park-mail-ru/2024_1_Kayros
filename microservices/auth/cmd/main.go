package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"2024_1_kayros/microservices/auth/internal/usecase"
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/user"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"2024_1_kayros/config"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.AuthGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice authorization doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice authorization responds on port %d", cfg.AuthGrpcServer.Port))

	// connecting to user microservice
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.AuthGrpcServer.Host, cfg.AuthGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("The microservice authorization is not available", zap.String("error", err.Error()))
	}
	defer func(userConn *grpc.ClientConn) {
		err := userConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(userConn)

	// init grpc server 
	server := grpc.NewServer()
	// register contract 
	client := user.NewUserManagerClient(userConn)
	auth.RegisterAuthManagerServer(server, usecase.NewLayer(client))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", cfg.AuthGrpcServer.Host, cfg.AuthGrpcServer.Port), zap.String("error", err.Error()))
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice authorization has shut down")
	os.Exit(0)
}
