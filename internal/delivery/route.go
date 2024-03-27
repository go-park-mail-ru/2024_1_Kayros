package delivery

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	deliveryAuth "2024_1_kayros/internal/delivery/auth"
	deliveryRest "2024_1_kayros/internal/delivery/restaurants"
	deliveryUser "2024_1_kayros/internal/delivery/user"
	"2024_1_kayros/internal/middleware"
	repoRest "2024_1_kayros/internal/repository/restaurants"
	repoSession "2024_1_kayros/internal/repository/session"
	repoUser "2024_1_kayros/internal/repository/user"
	usecaseAuth "2024_1_kayros/internal/usecase/auth"
	usecaseRest "2024_1_kayros/internal/usecase/restaurants"
	usecaseUser "2024_1_kayros/internal/usecase/user"
)

func Setup(db *sql.DB, redis *redis.Client, mux *mux.Router) {
	mux.PathPrefix("/api/v1")
	mux.StrictSlash(true)
	// слои repository
	rUser := repoUser.NewUserRepository(db)
	rSession := repoSession.NewSessionRepository(redis)
	rRestaurants := repoRest.NewRestaurantRepo(db)

	// слои usecase
	uAuth := usecaseAuth.NewAuthUsecase(rUser, rSession)
	uUser := usecaseUser.NewUserUsecase(rUser)
	uRestaurants := usecaseRest.NewRestaurantUseCase(rRestaurants)

	// слои delivery
	authHandlers := deliveryAuth.NewAuthDelivery(uAuth)
	userHandlers := deliveryUser.NewUserDelivery(uUser)
	restHandlers := deliveryRest.NewRestaurantHandler(uRestaurants)

	mux.HandleFunc("signin", authHandlers.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("signup", authHandlers.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("signout", authHandlers.SignOut).Methods("POST").Name("signout")

	mux.HandleFunc("user", userHandlers.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("restaurants", restHandlers.RestaurantList).Methods("GET").Name("restaurants")

	handler := middleware.SessionAuthentication(mux, rUser, rSession)
	handler = middleware.CorsMiddleware(handler)
}
