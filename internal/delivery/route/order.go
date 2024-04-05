package route

import (
	"database/sql"

	"github.com/gorilla/mux"

	"2024_1_kayros/internal/delivery/order"
	rOrder "2024_1_kayros/internal/repository/order"
	rUser "2024_1_kayros/internal/repository/user"
	ucOrder "2024_1_kayros/internal/usecase/order"
)

// нужно будет добавить интерфейс к БД и редису
func AddOrderRouter(db *sql.DB, mux *mux.Router) {
	repoOrder := rOrder.NewRepoLayer(db)
	repoUser := rUser.NewRepoLayer(db)
	usecaseOrder := ucOrder.NewUsecaseLayer(repoOrder, repoUser)
	handler := delivery.NewOrderHandler(usecaseOrder)

	mux.HandleFunc("order", handler.GetBasket).Methods("GET")
	mux.HandleFunc("order/update", handler.Update).Methods("PUT")
	mux.HandleFunc("order/update_status/{status}", handler.UpdateStatus).Methods("PUT")
	mux.HandleFunc("order/food/add/{food_id}", handler.AddFood).Methods("POST")
	mux.HandleFunc("order/food/update_count/", handler.UpdateFoodCount).Methods("PUT")
	mux.HandleFunc("order/food/delete/{food_id}", handler.DeleteFoodFromOrder).Methods("PUT")
}
