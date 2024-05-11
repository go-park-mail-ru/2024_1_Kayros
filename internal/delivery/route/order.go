package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	restproto "2024_1_kayros/gen/go/rest"
	userproto "2024_1_kayros/gen/go/user"
	delivery "2024_1_kayros/internal/delivery/order"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	ucOrder "2024_1_kayros/internal/usecase/order"
)

func AddOrderRouter(db *sql.DB, mux *mux.Router, userConn, restConn *grpc.ClientConn, logger *zap.Logger) {
	repoOrder := rOrder.NewRepoLayer(db)
	repoFood := rFood.NewRepoLayer(db)
	//init user grpc client
	grpcUserClient := userproto.NewUserManagerClient(userConn)
	//init rest grpc client
	grpcRestClient := restproto.NewRestWorkerClient(restConn)

	usecaseOrder := ucOrder.NewUsecaseLayer(repoOrder, repoFood, grpcUserClient, grpcRestClient)
	handler := delivery.NewOrderHandler(usecaseOrder, logger)

	mux.HandleFunc("/order", handler.GetBasket).Methods("GET")
	mux.HandleFunc("/order/{id}", handler.GetOrderById).Methods("GET")
	mux.HandleFunc("/orders/current", handler.GetCurrentOrders).Methods("GET")
	mux.HandleFunc("/order/update_address", handler.UpdateAddress).Methods("PUT")
	mux.HandleFunc("/order/pay", handler.Pay).Methods("PUT")
	mux.HandleFunc("/order/clean", handler.Clean).Methods("DELETE")
	mux.HandleFunc("/order/food/add", handler.AddFood).Methods("POST")
	mux.HandleFunc("/order/food/update_count", handler.UpdateFoodCount).Methods("PUT")
	mux.HandleFunc("/order/food/delete/{food_id}", handler.DeleteFoodFromOrder).Methods("DELETE")
}
