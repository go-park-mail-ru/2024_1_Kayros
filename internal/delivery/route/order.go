package route

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	restproto "2024_1_kayros/gen/go/rest"
	userproto "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	delivery "2024_1_kayros/internal/delivery/order"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	ucOrder "2024_1_kayros/internal/usecase/order"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"
)

func AddOrderRouter(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, logger *zap.Logger, metrics *metrics.Metrics) {
	repoOrder := rOrder.NewRepoLayer(cluster.PsqlClient, metrics)
	repoFood := rFood.NewRepoLayer(cluster.PsqlClient, metrics)
	//init user grpc client
	grpcUserClient := userproto.NewUserManagerClient(clients.UserConn)
	//init rest grpc client
	grpcRestClient := restproto.NewRestWorkerClient(clients.RestConn)

	usecaseOrder := ucOrder.NewUsecaseLayer(repoOrder, repoFood, grpcUserClient, grpcRestClient, metrics)
	handler := delivery.NewOrderHandler(usecaseOrder, logger)

	mux.HandleFunc("/api/v1/order", handler.GetBasket).Methods("GET")

	mux.HandleFunc("/api/v1/order/{id}", handler.GetOrderById).Methods("GET")

	mux.HandleFunc("/api/v1/promocode", handler.SetPromocode).Methods("POST")
	mux.HandleFunc("/api/v1/promocode", handler.GetAllPromocode).Methods("GET")

	mux.HandleFunc("/api/v1/orders/current", handler.GetCurrentOrders).Methods("GET")
	mux.HandleFunc("/api/v1/orders/archive", handler.GetArchiveOrders).Methods("GET")

	mux.HandleFunc("/api/v1/order/update_address", handler.UpdateAddress).Methods("PUT")
	mux.HandleFunc("/api/v1/order/pay", handler.Pay).Methods("PUT")
	mux.HandleFunc("/api/v1/order/clean", handler.Clean).Methods("DELETE")

	mux.HandleFunc("/api/v1/order/food/add", handler.AddFood).Methods("POST")
	mux.HandleFunc("/api/v1/order/food/update_count", handler.UpdateFoodCount).Methods("PUT")
	mux.HandleFunc("/api/v1/order/food/delete/{food_id}", handler.DeleteFoodFromOrder).Methods("DELETE")
}
