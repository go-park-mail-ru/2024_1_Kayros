package route

import (
	"database/sql"

	"github.com/gorilla/mux"

	"go.uber.org/zap"

	dRest "2024_1_kayros/internal/delivery/restaurants"
	rFood "2024_1_kayros/internal/repository/food"
	rRest "2024_1_kayros/internal/repository/restaurants"
	ucFood "2024_1_kayros/internal/usecase/food"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
)

func AddRestRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger) {
	repoRest := rRest.NewRepoLayer(db, logger)
	repoFood := rFood.NewRepoLayer(db, logger)
	usecaseRest := ucRest.NewUsecaseLayer(repoRest, logger)
	usecaseFood := ucFood.NewUsecaseLayer(repoFood, logger)
	deliveryRest := dRest.NewRestaurantHandler(usecaseRest, usecaseFood, logger)

	mux.HandleFunc("/restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("/restaurants/{id}", deliveryRest.RestaurantById).Methods("GET").Name("restaurants-detail")
}
