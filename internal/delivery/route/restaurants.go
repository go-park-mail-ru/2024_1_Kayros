package route

import (
	"database/sql"

	dRest "2024_1_kayros/internal/delivery/restaurants"
	rRest "2024_1_kayros/internal/repository/restaurants"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// нужно будет добавить интерфейс к БД и редису
func AddRestRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger) {
	repoRest := rRest.NewRepoLayer(db, logger)
	usecaseRest := ucRest.NewUsecaseLayer(repoRest, logger)
	deliveryRest := dRest.NewDelivery(usecaseRest, logger)

	mux.HandleFunc("restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("restaurants/{id}", deliveryRest.RestaurantById).Methods("GET").Name("restaurants-detail")
}
