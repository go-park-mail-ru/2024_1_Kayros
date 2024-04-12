package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/route"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/minio"
	"2024_1_kayros/services/postgres"
	"2024_1_kayros/services/redis"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Run создает все сервисы приложения и запускает их
func Run(cfg *config.Project) {
	logger := zap.Must(zap.NewProduction())
	functions.InitValidator(logger)

	postgreDB := postgres.Init(cfg, logger)
	redisSessionDB := redis.Init(cfg, logger, cfg.Redis.DatabaseSession)
	redisCsrfDB := redis.Init(cfg, logger, cfg.Redis.DatabaseCsrf)
	minioDB := minio.Init(cfg, logger)

	r := mux.NewRouter()
	handler := route.Setup(cfg, postgreDB, redisSessionDB, redisCsrfDB, minioDB, r, logger)

	serverConfig := cfg.Server
	serverAddress := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Handler:      handler,
		Addr:         serverAddress,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second, // таймаут на запись данных в ответ на запрос
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,  // таймаут на чтение данных из запроса
		IdleTimeout:  time.Duration(serverConfig.IdleTimeout) * time.Second,  // время поддержания связи между клиентом и сервером
	}
	//log.Fatal(srv.ListenAndServe())
	srvConfig := cfg.Server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Сервер не может быть запущен.\n%v", err)
		} else {
			log.Printf("Сервер запущен по адресу %s:%d", srvConfig.Host, srvConfig.Port)
		}
	}()

	// канал для получения прерывания, завершающего работу сервиса (ожидает прерывание процесса)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(srvConfig.ShutdownDuration))
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Сервер экстренно завершил свою работу с ошибкой.\n%v", err)
		os.Exit(1) //
	}

	log.Println("Сервер завершил свою работу")
	os.Exit(0)
}
