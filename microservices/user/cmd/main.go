package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	metrics "2024_1_kayros/microservices/metrics"
	"2024_1_kayros/microservices/user/internal/repo"
	"2024_1_kayros/microservices/user/internal/repo/stmts"
	"2024_1_kayros/microservices/user/internal/usecase"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/user"
	grpcServerMiddleware "2024_1_kayros/internal/middleware/grpc/server"
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
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg, "user")
	middleware := grpcServerMiddleware.NewMiddlewareChain(logger, metrics)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		address := fmt.Sprintf("%s:%d", cfg.UserGrpcServerExporter.Host, cfg.UserGrpcServerExporter.Port)
		logger.Info(fmt.Sprintf("Serving metrics responds on port %d", cfg.UserGrpcServerExporter.Port))
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.Fatal("Error starting metrics server", zap.String("error", err.Error()))
		}
	}()

	// init grpc server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.MetricsMiddleware, middleware.AccessMiddleware))
	// init services for server work
	postgreDB := postgres.Init(cfg, logger)
	minioClient := minio.Init(cfg, logger)
	// init prepared statements
	statements, err := stmts.InitPrepareStatements(postgreDB)
	defer func (stmts map[string]*sql.Stmt) {
		for _, stmt := range statements {
			stmt.Close()
		}
	}(statements)
	if err != nil {
		logger.Fatal("Can't define prepared statements for user database")
	}
    // init repo layers
	repoUser := repo.NewLayer(postgreDB, metrics, statements)
	repoMinio := minios3.NewRepoLayer(minioClient)
	user.RegisterUserManagerServer(server, usecase.NewLayer(repoUser, repoMinio, logger))
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
