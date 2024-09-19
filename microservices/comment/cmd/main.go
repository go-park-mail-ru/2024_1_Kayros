package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	cfg "2024_1_kayros/config"
	"2024_1_kayros/gen/go/comment"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
	"2024_1_kayros/microservices/comment/internal/repo"
	"2024_1_kayros/microservices/comment/internal/usecase"
	metrics "2024_1_kayros/microservices/metrics"
	"2024_1_kayros/services/postgres"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg.Read(logger)
	projConfig := cfg.Config

	port := fmt.Sprintf(":%d", projConfig.CommentGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice comment doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice comment responds on port %d", projConfig.CommentGrpcServer.Port))
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "comment")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf("%s:%d", projConfig.CommentGrpcServerExporter.Host, projConfig.CommentGrpcServerExporter.Port)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", projConfig.CommentGrpcServerExporter.Port))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	//init grpc server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	//init services for server work
	postgreDB := postgres.Init(logger)
	repoComment := repo.NewCommentLayer(postgreDB, metrics)
	comment.RegisterCommentWorkerServer(server, usecase.NewCommentLayer(repoComment, logger))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", projConfig.CommentGrpcServer.Host, projConfig.CommentGrpcServer.Port), zap.String("error", err.Error()))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice comment has shut down")
	os.Exit(0)
}
