package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/rest"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
	metrics "2024_1_kayros/microservices/metrics"
	"2024_1_kayros/microservices/restaurants/internal/repo"
	"2024_1_kayros/microservices/restaurants/internal/usecase"
	"2024_1_kayros/services/postgres"
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
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "rest")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)

	// init grpc server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	// init services for server work
	postgreDB := postgres.Init(cfg, logger)
	repoRest := repo.NewRestLayer(postgreDB, metrics)
	rest.RegisterRestWorkerServer(server, usecase.NewRestLayer(repoRest, logger))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", cfg.RestGrpcServer.Host, cfg.RestGrpcServer.Port), zap.String("error", err.Error()))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice restaurant has shut down")
	os.Exit(0)
}
