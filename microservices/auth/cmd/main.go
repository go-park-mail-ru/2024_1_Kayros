package main

import (
	"fmt"
	"log"
	"net"

	"2024_1_kayros/microservices/auth/internal/usecase"
	authv1 "2024_1_kayros/microservices/auth/proto"
	userv1 "2024_1_kayros/microservices/user/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/functions"
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
	infoMsg := fmt.Sprintf("The authentication server listens port %d", cfg.AuthGrpcServer.Port)
	logger.Info(infoMsg)

	// init grpc server
	server := grpc.NewServer()

	// connect to user microservice
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.UserGrpcServer.Host, cfg.UserGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := fmt.Sprintf("The user microservice is not available.\n%v", err)
		logger.Error(errMsg)
	}
	defer func(userConn *grpc.ClientConn) {
		err := userConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(userConn)
	client := userv1.NewUserManagerClient(userConn)

	authv1.RegisterAuthManagerServer(server, usecase.NewLayer(client))
	err = server.Serve(conn)
	if err != nil {
		log.Fatalf("error in serving server on port %d -  %s", cfg.CommentGrpcServer.Port, err)
	}
	
	
}
