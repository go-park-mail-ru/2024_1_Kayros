package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"2024_1_kayros/internal/delivery/order"
	rOrder "2024_1_kayros/internal/repository/order"
	rUser "2024_1_kayros/internal/repository/user"
	ucOrder "2024_1_kayros/internal/usecase/order"
)

// нужно будет добавить интерфейс к БД и редису
func AddOrderRouter(db *sql.DB, minio *minio.Client, mux *mux.Router, logger *zap.Logger) {
	repoOrder := rOrder.NewRepoLayer(db, logger)
	repoUser := rUser.NewRepoLayer(db, minio, logger)
	usecaseOrder := ucOrder.NewUsecaseLayer(repoOrder, repoUser, logger)
	handler := delivery.NewOrderHandler(usecaseOrder, logger)

	mux.HandleFunc("/order", handler.GetBasket).Methods("GET") //работает
	mux.HandleFunc("/order/update_address", handler.UpdateAddress).Methods("PUT")
	mux.HandleFunc("/order/pay", handler.Pay).Methods("PUT")
	mux.HandleFunc("/order/food/add/{food_id}", handler.AddFood).Methods("POST")                  //работает
	mux.HandleFunc("/order/food/update_count", handler.UpdateFoodCount).Methods("PUT")            //работает
	mux.HandleFunc("/order/food/delete/{food_id}", handler.DeleteFoodFromOrder).Methods("DELETE") //работает
}
