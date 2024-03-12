package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"2024_1_kayros/internal/delivery/authorization"
	"2024_1_kayros/internal/delivery/middlewares"
	"2024_1_kayros/internal/entity"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(false)

	// флаг для установки времени graceful shutdown-а
	var wait time.Duration
	flag.DurationVar(&wait, "grtm", time.Second*15, "Промежуток времени, в течение которого сервер "+
		"плавно завершает работу, завершая текущие запросы")
	flag.Parse()

	// для работы блока авторизации необходима база данных
	auth := authorization.AuthHandler{
		Database: entity.InitDatabase(),
	}

	// авторизация, регистрация, деавторизация
	r.HandleFunc("/api/v1/signin/", auth.SignIn).Methods("POST").Name("signin")
	r.HandleFunc("/api/v1/signup/", auth.SignUp).Methods("POST").Name("signup")
	r.HandleFunc("/api/v1/signout/", auth.SignOut).Methods("POST").Name("signout")

	// рестораны
	r.HandleFunc("api/v1/restaurants/", auth.Database.Restaurants.RestaurantList).Methods("GET").Name("restaurants")

	// устанавливаем middlewares для аутентификации с помощью сессионной куки
	handler := middlewares.SessionAuthentication(r, &auth.Database)
	// устанавливаем middleware для CORS
	handler = middlewares.CorsMiddleware(handler)
	//srv := &http.Server{
	//	Handler:      r,
	//	Addr:         ":8000",
	//	WriteTimeout: 10 * time.Second, // таймаут на запись данных в ответ на запрос
	//	ReadTimeout:  10 * time.Second, // таймаут на чтение данных из запроса
	//	IdleTimeout:  30 * time.Second, // время поддержания связи между клиентом и сервером
	//}
	log.Fatal(http.ListenAndServe(":8000", r))
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
