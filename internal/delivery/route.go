package delivery

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	deliveryAuth "2024_1_kayros/internal/delivery/auth"
	deliveryOrder "2024_1_kayros/internal/delivery/order"
	deliveryRest "2024_1_kayros/internal/delivery/restaurants"
	deliveryUser "2024_1_kayros/internal/delivery/user"
	"2024_1_kayros/internal/middleware"
	repoFood "2024_1_kayros/internal/repository/food"
	repoOrder "2024_1_kayros/internal/repository/order"
	repoRest "2024_1_kayros/internal/repository/restaurants"
	repoSession "2024_1_kayros/internal/repository/session"
	repoUser "2024_1_kayros/internal/repository/user"
	usecaseAuth "2024_1_kayros/internal/usecase/auth"
	usecaseFood "2024_1_kayros/internal/usecase/food"
	usecaseOrder "2024_1_kayros/internal/usecase/order"
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
	rFood := repoFood.NewFoodRepository(db)
	rOrder := repoOrder.NewRepository(db)

	// слои usecase
	uAuth := usecaseAuth.NewAuthUsecase(rUser, rSession)
	uUser := usecaseUser.NewUserUsecase(rUser)
	uRestaurants := usecaseRest.NewRestaurantUseCase(rRestaurants)
	uFood := usecaseFood.NewUseCase(rFood)
	uOrder := usecaseOrder.NewUseCase(rOrder, rUser)

	// слои delivery
	authHandlers := deliveryAuth.NewAuthDelivery(uAuth)
	userHandlers := deliveryUser.NewUserDelivery(uUser)
	restHandlers := deliveryRest.NewRestaurantHandler(uRestaurants, uFood)
	orderHandlers := deliveryOrder.NewOrderHandler(uOrder)

	mux.HandleFunc("signin", authHandlers.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("signup", authHandlers.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("signout", authHandlers.SignOut).Methods("POST").Name("signout")

	mux.HandleFunc("user", userHandlers.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("restaurants", restHandlers.RestaurantList).Methods("GET").Name("restaurants")
	mux.HandleFunc("restaurants/{id}", restHandlers.RestaurantById).Methods("GET").Name("restaurant")

	mux.HandleFunc("order", orderHandlers.GetOrder).Methods("GET")
	mux.HandleFunc("order/update", orderHandlers.Update).Methods("PUT")
	mux.HandleFunc("order/food/add/{food_id}", orderHandlers.AddFood).Methods("POST")
	mux.HandleFunc("order/food/update_count/", orderHandlers.UpdateFoodCount).Methods("PUT")
	mux.HandleFunc("order/food/delete/{food_id}", orderHandlers.DeleteFoodFromOrder).Methods("PUT")

	handler := middleware.SessionAuthentication(mux, rUser, rSession)
	handler = middleware.CorsMiddleware(handler)
}
