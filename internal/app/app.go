package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"2024_1_kayros/config"
	route "2024_1_kayros/internal/delivery"
	"github.com/gorilla/mux"
)

// Run создает все сервисы приложения и запускает их
func Run(cfg *config.Config) {
	// вот тут вот нужно создать редис, постгрес, минио
	// ....
	r := mux.NewRouter()
	r.StrictSlash(true)

	// нужно будет поменять на настоящий объект базы данных
	route.Setup(&sql.DB{}, r)

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

	// канал для получения прерывания, завершающего работу сервиса (ожидает Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	// контекст ожидания выполнения запросов в течение времени wait
	ctx, cancel := context.WithTimeout(context.Background(), srvConfig.ShutdownDuration)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Сервер завершил свою работу с ошибкой.\n%v", err)
		os.Exit(1) //
	}

	log.Println("Сервер завершил свою работу успешно")
	os.Exit(0)
}
