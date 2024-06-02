package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/rest"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
	metrics "2024_1_kayros/microservices/metrics"
	"2024_1_kayros/microservices/restaurants/internal/repo"
	"2024_1_kayros/microservices/restaurants/internal/repo/stmts"
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
	metrics := metrics.NewMetrics(reg, "restaurants")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf("%s:%d", cfg.RestGrpcServerExporter.Host, cfg.RestGrpcServerExporter.Port)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.RestGrpcServerExporter.Port))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	// init grpc server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	// init services for server work
	postgreDB := postgres.Init(cfg, logger)
	// init prepared statements
	
	statements, err := stmts.InitPrepareStatements(postgreDB)
	if err != nil {
		logger.Fatal("Can't define prepared statements for user database")
	}
	defer func (stmts map[string]*sql.Stmt) {
		for _, stmt := range statements {
			stmt.Close()
		}
	}(statements)
	repoRest := repo.NewRestLayer(postgreDB, metrics, statements)
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
