package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	cfg "2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/delivery/route"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Run creates all services and run the handler goroutines
func Run(logger *zap.Logger) {
	// add custom struct tags for dto validation
	functions.InitDtoValidator(logger)
	projectCfg := cfg.Config
	// initialization of cluster services
	serviceCluster := services.Init(logger)
	// initialization of metrics client
	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg, "gateway")
	// initialization of grpc clients
	grpcClients := microservices.Init(logger, m)
	defer func(clients *microservices.Clients) {
		err := clients.RestConn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while closing connection with microservice 'restaurant': %v", err))
		}
		err = clients.CommentConn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while closing connection with microservice 'comment': %v", err))
		}
		err = clients.AuthConn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while closing connection with microservice 'authorization': %v", err))
		}
		err = clients.UserConn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while closing connection with microservice 'user': %v", err))
		}
		err = clients.SessionConn.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while closing connection with microservice 'session': %v", err))
		}
	}(grpcClients)

	r := mux.NewRouter()
	handler := route.Setup(r, serviceCluster, grpcClients, m, reg, logger)

	serverConfig := projectCfg.Server
	serverAddress := fmt.Sprintf(":%d", serverConfig.Port)
	srv := &http.Server{
		Handler:      handler,
		Addr:         serverAddress,
		WriteTimeout: serverConfig.WriteTimeout, // timeout for writing data in response to a request
		ReadTimeout:  serverConfig.ReadTimeout,  // timeout for reading data from request
		IdleTimeout:  serverConfig.IdleTimeout,  // time of communication between client and server
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("the main service has finished with error: %v", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), serverConfig.ShutdownDuration)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("the main service urgently shut down with an error: %v", err))
	}
	logger.Info("the main service has shut down")
}
