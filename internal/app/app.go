package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	metrics "2024_1_kayros"
	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/route"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/minio"
	"2024_1_kayros/services/postgres"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Run creates all services and run the handler goroutines
func Run(cfg *config.Project) {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)

	postgreDB := postgres.Init(cfg, logger)
	minioDB := minio.Init(cfg, logger)
	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)

	//restaurant microservice
	restConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.RestGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := fmt.Sprintf("The microservice restaurant is not available.\n%v", err)
		logger.Error(errMsg)
	}
	defer func(restConn *grpc.ClientConn) {
		err := restConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(restConn)

	//comment microservice
	commentConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.CommentGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := fmt.Sprintf("The microservice comment is not available.\n%v", err)
		logger.Error(errMsg)
	}
	defer func(restConn *grpc.ClientConn) {
		err := restConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(commentConn)

	//auth microservice
	authConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.AuthGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := fmt.Sprintf("The auth microservice is not available.\n%v", err)
		logger.Error(errMsg)
	}
	defer func(authConn *grpc.ClientConn) {
		err := authConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(authConn)

	// user microservice
	userConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.UserGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := fmt.Sprintf("The microservice user is not available.\n%v", err)
		logger.Error(errMsg)
	}
	defer func(userConn *grpc.ClientConn) {
		err := userConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(userConn)

	// session microservice
	sessionConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.SessionGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := fmt.Sprintf("The microservice session is not available.\n%v", err)
		logger.Error(errMsg)
	}
	defer func(sessionConn *grpc.ClientConn) {
		err := userConn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(sessionConn)

	r := mux.NewRouter()
	handler := route.Setup(cfg, postgreDB, minioDB, r, logger, restConn, commentConn, authConn, userConn, sessionConn, m)

	serverConfig := cfg.Server
	serverAddress := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Handler:      handler,
		Addr:         serverAddress,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second, // timeout for writing data in response to a request
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,  // timeout for reading data from request
		IdleTimeout:  time.Duration(serverConfig.IdleTimeout) * time.Second,  // time of communication between client and server
	}

	srvConfig := cfg.Server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("The main service cannot be started.\n%v", err)
		} else {
			log.Printf("The main service has started at the address %s:%d", srvConfig.Host, srvConfig.Port)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(srvConfig.ShutdownDuration))
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Printf("The main service urgently shut down with an error.\n%v", err)
		os.Exit(1) //
	}

	log.Println("The main service has shut down")
	os.Exit(0)
}
