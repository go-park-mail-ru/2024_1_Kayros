package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"2024_1_kayros/internal/delivery/authorization"
	"2024_1_kayros/internal/delivery/middlewares"
	"2024_1_kayros/internal/delivery/restaurants"
	"2024_1_kayros/internal/entity"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	const PORT int = 8000 // нужно в файл конфигурации закинуть

	// флаг для установки времени graceful shutdown-а
	var wait time.Duration
	flag.DurationVar(&wait, "grtm", time.Second*15, "Промежуток времени, в течение которого сервер "+
		"плавно завершает работу, завершая текущие запросы")
	flag.Parse()

	// все в файл конфигурации
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"

	// Initialize minio client object.
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called testbucket.
	bucketName := "testbucket"
	location := "us-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "testdata"
	filePath := "/tmp/testdata"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	// для работы блока авторизации
	auth := authorization.AuthHandler{
		DB: entity.InitDatabase(),
	}
	// для работы с ресторанами
	rest := restaurants.InitRestaurantStore()

	// авторизация, регистрация, деавторизация
	r.HandleFunc("/docs/v1/signin", auth.SignIn).Methods("POST").Name("signin")
	r.HandleFunc("/docs/v1/signup", auth.SignUp).Methods("POST").Name("signup")
	r.HandleFunc("/docs/v1/signout", auth.SignOut).Methods("POST").Name("signout")
	// получение информации о пользователе
	r.HandleFunc("/docs/v1/user", auth.UserData).Methods("GET").Name("userdata")
	// рестораны
	r.HandleFunc("/docs/v1/restaurants", rest.RestaurantList).Methods("GET").Name("restaurants")

	// устанавливаем middlewares для аутентификации с помощью сессионной куки
	handler := middlewares.SessionAuthentication(r, &auth.DB)
	// устанавливаем middleware для CORS
	handler = middlewares.CorsMiddleware(handler)
	//srv := &http.Server{
	//	Handler:      r,
	//	Addr:         ":8000",
	//	WriteTimeout: 10 * time.Second, // таймаут на запись данных в ответ на запрос
	//	ReadTimeout:  10 * time.Second, // таймаут на чтение данных из запроса
	//	IdleTimeout:  30 * time.Second, // время поддержания связи между клиентом и сервером
	//}
	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), handler))
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
