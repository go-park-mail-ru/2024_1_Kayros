package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/internal/delivery/order"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	rUser "2024_1_kayros/internal/repository/user"
	ucOrder "2024_1_kayros/internal/usecase/order"
)

func AddOrderRouter(db *sql.DB, mux *mux.Router, logger *zap.Logger) {
	repoOrder := rOrder.NewRepoLayer(db, logger)
	repoFood := rFood.NewRepoLayer(db)
	repoUser := rUser.NewRepoLayer(db)
	usecaseOrder := ucOrder.NewUsecaseLayer(repoOrder, repoFood, repoUser, logger)
	handler := delivery.NewOrderHandler(usecaseOrder, logger)

	mux.HandleFunc("/order", handler.GetBasket).Methods("GET")
	mux.HandleFunc("/order/current", handler.GetOrders).Methods("GET")
	mux.HandleFunc("/order/update_address", handler.UpdateAddress).Methods("PUT")
	mux.HandleFunc("/order/pay", handler.Pay).Methods("PUT")
	mux.HandleFunc("/order/clean", handler.Clean).Methods("DELETE")
	mux.HandleFunc("/order/food/add", handler.AddFood).Methods("POST")
	mux.HandleFunc("/order/food/update_count", handler.UpdateFoodCount).Methods("PUT")
	mux.HandleFunc("/order/food/delete/{food_id}", handler.DeleteFoodFromOrder).Methods("DELETE")
}
