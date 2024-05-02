package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/route"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/minio"
	"2024_1_kayros/services/postgres"
	"2024_1_kayros/services/redis"
)

// Run creates all services and run the handler goroutines
func Run(cfg *config.Project) {
	logger := zap.Must(zap.NewProduction())
	functions.InitDtoValidator(logger)

	postgreDB := postgres.Init(cfg, logger)
	redisSessionDB := redis.Init(cfg, logger, cfg.Redis.DatabaseSession)
	redisCsrfDB := redis.Init(cfg, logger, cfg.Redis.DatabaseCsrf)
	//redisUnauthTokensDB := redis.Init(cfg, logger, cfg.Redis.DatabaseUnauthTokens)
	minioDB := minio.Init(cfg, logger)

	//rest microservice
	restConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.RestGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("The server cannot be started.\n%v", err)
	} else {
		log.Printf("The server is started at the address %s:%d", cfg.RestGrpcServer.Host, cfg.RestGrpcServer.Port)
	}
	defer restConn.Close()

	//comment microservice
	commentConn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.CommentGrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("The server cannot be started.\n%v", err)
	} else {
		log.Printf("The server is started at the address %s:%d", cfg.CommentGrpcServer.Host, cfg.CommentGrpcServer.Port)
	}
	defer restConn.Close()

	r := mux.NewRouter()
	handler := route.Setup(cfg, postgreDB, redisSessionDB, redisCsrfDB, minioDB, r, logger, restConn, commentConn)

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
			log.Printf("The server cannot be started.\n%v", err)
		} else {
			log.Printf("The server is started at the address %s:%d", srvConfig.Host, srvConfig.Port)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(srvConfig.ShutdownDuration))
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Printf("The server urgently shut down with an error.\n%v", err)
		os.Exit(1) //
	}

	log.Println("The server has shut down")
	os.Exit(0)
}
