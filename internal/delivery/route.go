package delivery

import (
	"database/sql"

	authDelivery "2024_1_kayros/internal/delivery/auth"
	"2024_1_kayros/internal/delivery/restaurants"
	"2024_1_kayros/internal/delivery/user"
	"2024_1_kayros/internal/middleware"
	repoSession "2024_1_kayros/internal/repository/session"
	repoUser "2024_1_kayros/internal/repository/user"
	authUsecase "2024_1_kayros/internal/usecase/auth"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func Setup(db *sql.DB, redis *redis.Client, mux *mux.Router) {
	userRepo := repoUser.NewUserRepository(db)
	sessionRepo := repoSession.NewSessionRepository(redis)

	userUsecase := authUsecase.NewAuthUsecase(userRepo, sessionRepo)

	authHandlers := authDelivery.NewAuthDelivery(userUsecase)

	mux.HandleFunc("/api/v1/signin", authHandlers.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/api/v1/signup", authHandlers.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/api/v1/signout", authHandlers.SignOut).Methods("POST").Name("signout")

	mux.HandleFunc("/api/v1/user", user.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("/api/v1/restaurants", restaurants.RestaurantList).Methods("GET").Name("restaurants")

	handler := middleware.SessionAuthentication(mux)
	handler = middleware.CorsMiddleware(handler)
}
