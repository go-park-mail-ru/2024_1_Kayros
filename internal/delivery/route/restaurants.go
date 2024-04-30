package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"go.uber.org/zap"

	dRest "2024_1_kayros/internal/delivery/restaurants"
	rFood "2024_1_kayros/internal/repository/food"
	ucFood "2024_1_kayros/internal/usecase/food"
	ucRest "2024_1_kayros/internal/usecase/restaurants"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

func AddRestRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger, conn *grpc.ClientConn) {
	//repoRest := rRest.NewRepoLayer(db)
	repoFood := rFood.NewRepoLayer(db)
	usecaseFood := ucFood.NewUsecaseLayer(repoFood)

	GrpcRest := rest.NewRestWorkerClient(conn)
	usecaseRest := ucRest.NewUsecaseLayer(GrpcRest)

	deliveryRest := dRest.NewRestaurantHandler(usecaseRest, usecaseFood, logger)

	mux.HandleFunc("/restaurants", deliveryRest.RestaurantList).Methods("GET").Name("restaurants-list")
	mux.HandleFunc("/restaurants/{id}", deliveryRest.RestaurantById).Methods("GET").Name("restaurants-detail")
}
