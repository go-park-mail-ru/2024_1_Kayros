package route

import (
	"database/sql"

	"github.com/gorilla/mux"

	"2024_1_kayros/internal/delivery/restaurants"
	rFood "2024_1_kayros/internal/repository/food"
	rRest "2024_1_kayros/internal/repository/restaurants"
	ucFood "2024_1_kayros/internal/usecase/food"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
)

// нужно будет добавить интерфейс к БД и редису
func AddRestRouter(db *sql.DB, mux *mux.Router) {
	repoRest := rRest.NewRepoLayer(db)
	repoFood := rFood.NewRepoLayer(db)
	usecaseFood := ucFood.NewUsecaseLayer(repoFood)
	usecaseRest := ucRest.NewUsecaseLayer(repoRest)
	deliveryRest := delivery.NewRestaurantHandler(usecaseRest, usecaseFood)

	mux.HandleFunc("restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("restaurants/{id}", deliveryRest.RestaurantById).Methods("GET").Name("restaurants-detail")
}
