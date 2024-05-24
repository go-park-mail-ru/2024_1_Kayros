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
	"2024_1_kayros/gen/go/comment"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
	"2024_1_kayros/microservices/comment/internal/repo"
	"2024_1_kayros/microservices/comment/internal/usecase"
	metrics "2024_1_kayros/microservices/metrics"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg := config.NewConfig(logger)

	port := fmt.Sprintf(":%d", cfg.CommentGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice comment doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice comment responds on port %d", cfg.CommentGrpcServer.Port))
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "food")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)


	//init grpc server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	//init services for server work
	postgreDB := postgres.Init(cfg, logger)
	repoComment := repo.NewCommentLayer(postgreDB, metrics)
	comment.RegisterCommentWorkerServer(server, usecase.NewCommentLayer(repoComment, logger))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", cfg.CommentGrpcServer.Host, cfg.CommentGrpcServer.Port), zap.String("error", err.Error()))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice comment has shut down")
	os.Exit(0)
}
