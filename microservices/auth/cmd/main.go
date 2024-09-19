package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/user"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
	"2024_1_kayros/microservices/auth/internal/usecase"
	metrics "2024_1_kayros/microservices/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cfg "2024_1_kayros/config"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg.Read(logger)
	projConfig := cfg.Config

	port := fmt.Sprintf(":%d", projConfig.AuthGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice authorization doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice authorization responds on port %d", projConfig.AuthGrpcServer.Port))
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "auth")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf("%s:%d", projConfig.AuthGrpcServerExporter.Host, projConfig.AuthGrpcServerExporter.Port)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", projConfig.AuthGrpcServerExporter.Port))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	// connecting to user microservice
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", projConfig.UserGrpcServer.Host, projConfig.UserGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	// register contract
	client := user.NewUserManagerClient(userConn)
	auth.RegisterAuthManagerServer(server, usecase.NewLayer(client, logger))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", projConfig.AuthGrpcServer.Host, projConfig.AuthGrpcServer.Port), zap.String("error", err.Error()))
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice authorization has shut down")
	os.Exit(0)
}
