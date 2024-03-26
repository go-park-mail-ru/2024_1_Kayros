package delivery

import (
	"database/sql"

	deliveryAuth "2024_1_kayros/internal/delivery/auth"
	deliveryRest "2024_1_kayros/internal/delivery/restaurants"
	deliveryUser "2024_1_kayros/internal/delivery/user"
	"2024_1_kayros/internal/middleware"
	repoSession "2024_1_kayros/internal/repository/session"
	repoUser "2024_1_kayros/internal/repository/user"
	usecaseAuth "2024_1_kayros/internal/usecase/auth"
	usecaseUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func Setup(db *sql.DB, redis *redis.Client, mux *mux.Router) {
	mux.PathPrefix("/api/v1")
	mux.StrictSlash(true)
	// слои repository
	rUser := repoUser.NewUserRepository(db)
	rSession := repoSession.NewSessionRepository(redis)

	// слои usecase
	uAuth := usecaseAuth.NewAuthUsecase(rUser, rSession)
	uUser := usecaseUser.NewUserUsecase(rUser)

	// слои delivery
	authHandlers := deliveryAuth.NewAuthDelivery(uAuth)
	userHandlers := deliveryUser.NewUserDelivery(uUser)

	mux.HandleFunc("signin", authHandlers.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("signup", authHandlers.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("signout", authHandlers.SignOut).Methods("POST").Name("signout")

	mux.HandleFunc("user", userHandlers.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants")

	handler := middleware.SessionAuthentication(mux, rUser, rSession)
	handler = middleware.CorsMiddleware(handler)
}
