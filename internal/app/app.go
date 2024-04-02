package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/route"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/postgres"
	"2024_1_kayros/services/redis"
	"github.com/gorilla/mux"
)

// Run создает все сервисы приложения и запускает их
func Run(cfg *config.Project) {
	functions.InitValidator()
	// вот тут вот нужно создать редис, минио
	// ....
	postgreDB, err := postgres.Init(cfg)
	if err != nil {
		log.Printf("Не удалось подключиться к базе данных %s по адресу %s:%d\n%s\n",
			cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port, err)
		return
	}

	redisDB, err := redis.Init(cfg)

	r := mux.NewRouter()
	route.Setup(postgreDB, redisDB, r)

	srvConfig := cfg.Server
	srvAddress := srvConfig.Host + ":" + strconv.Itoa(srvConfig.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         srvAddress,
		WriteTimeout: srvConfig.WriteTimeout, // таймаут на запись данных в ответ на запрос
		ReadTimeout:  srvConfig.ReadTimeout,  // таймаут на чтение данных из запроса
		IdleTimeout:  srvConfig.IdleTimeout,  // время поддержания связи между клиентом и сервером
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), srvConfig.ShutdownDuration)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Сервер экстренно завершил свою работу с ошибкой.\n%v", err)
		os.Exit(1) //
	}

	log.Println("Сервер завершил свою работу успешно")
	os.Exit(0)
}
