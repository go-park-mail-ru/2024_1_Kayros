package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"2024_1_kayros/internal/delivery/authorization"
	"2024_1_kayros/internal/delivery/restaurants"
)

func main() {
	r := mux.NewRouter()

	// мультиплексор авторизации
	auth := authorization.NewAuthStore()
	restaurants := delivery.NewRestaurantStore()

	// флаг для установки времени graceful shutdown-а
	var wait time.Duration
	flag.DurationVar(&wait, "grtm", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// устанавливаем middleware для аутентификации с помощью сессионной куки
	r.Use(auth.SessionAuthentication)

	// cаброутинг для блока авторизации
	subRoutingAuth := r.Methods("POST").Subrouter()
	// авторизация, регистрация, деавторизация
	subRoutingAuth.HandleFunc("/signin", auth.SignIn).Name("signin")
	subRoutingAuth.HandleFunc("/signup", auth.SignUp).Name("signup")
	subRoutingAuth.HandleFunc("/signout", auth.SignOut).Name("signout")

	// рестораны
	r.HandleFunc("/restaurants", restaurants.RestaurantList).Methods("GET").Name("restaurants")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 10 * time.Second, // таймаут на запись данных в ответ на запрос
		ReadTimeout:  10 * time.Second, // таймаут на чтение данных из запроса
		IdleTimeout:  30 * time.Second, // время поддержания связи между клиентом и сервером
	}
	log.Fatal(srv.ListenAndServe())
	//
	//go func() {
	//	if err := srv.ListenAndServe(); err != nil {
	//		log.Printf("Сервер не может быть запущен. Ошибка: \n%v", err)
	//	} else {
	//		log.Println("Сервер запущен на порте 8000")
	//	}
	//}()
	//
	//// канал для получения прерывания, завершающего работу сервиса (ожидает Ctrl+C)
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//
	//<-c
	//
	//// контекст ожидания выполнения запросов в течение времени wait
	//ctx, cancel := context.WithTimeout(context.Background(), wait)
	//defer cancel()
	//err := srv.Shutdown(ctx)
	//if err != nil {
	//	log.Printf("Сервер завершил свою работу с ошибкой. Ошибка: \n%v", err)
	//	os.Exit(1) //
	//}
	//
	//log.Println("Сервер завершил свою работу успешно")
	//os.Exit(0)
}
