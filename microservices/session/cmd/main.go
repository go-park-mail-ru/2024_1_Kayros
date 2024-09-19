package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"2024_1_kayros/gen/go/session"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
	metrics "2024_1_kayros/microservices/metrics"
	"2024_1_kayros/microservices/session/internal/repo"
	"2024_1_kayros/microservices/session/internal/usecase"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	cfg "2024_1_kayros/config"
	"2024_1_kayros/services/redis"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	cfg.Read(logger)
	projectCfg := cfg.Config

	port := fmt.Sprintf(":%d", projectCfg.SessionGrpcServer.Port)
	conn, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("The microservice session doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("The microservice session responds on port %d", projectCfg.SessionGrpcServer.Port))
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "session")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf("%s:%d", projectCfg.SessionGrpcServerExporter.Host, projectCfg.SessionGrpcServerExporter.Port)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", projectCfg.SessionGrpcServerExporter.Port))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	redisSession := redis.Init(&projectCfg, logger, projectCfg.Redis.DatabaseSession)
	redisCsrf := redis.Init(&projectCfg, logger, projectCfg.Redis.DatabaseCsrf)

	repoSession := repo.NewLayer(redisSession, metrics)
	repoCsrf := repo.NewLayer(redisCsrf, metrics)

	session.RegisterSessionManagerServer(server, usecase.NewLayer(repoCsrf, repoSession, logger, &projectCfg.Redis))
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error serving on %s:%d", projectCfg.SessionGrpcServer.Host, projectCfg.SessionGrpcServer.Port), zap.String("error", err.Error()))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("The microservice session has shut down")
	os.Exit(0)
}
