package route

import (
	restproto "2024_1_kayros/gen/go/rest"
	sessionproto "2024_1_kayros/gen/go/session"
	userproto "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/delivery/payment"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/usecase/order"
	ucSession "2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

func AddPaymentRouter(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, logger *zap.Logger, metrics *metrics.Metrics) {
	repoFood := rFood.NewRepoLayer(cluster.PsqlClient, metrics)
	repoOrder := rOrder.NewRepoLayer(cluster.PsqlClient, metrics)

	// init session grpc client
	grpcSessionClient := sessionproto.NewSessionManagerClient(clients.SessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient, metrics)
	// init user grpc client
	grpcUserClient := userproto.NewUserManagerClient(clients.UserConn)
	// init rest grpc client
	grpcRestClient := restproto.NewRestWorkerClient(clients.RestConn)

	usecaseOrder := order.NewUsecaseLayer(repoOrder, repoFood, grpcUserClient, grpcRestClient, metrics)
	deliveryPayment := payment.NewPaymentDelivery(usecaseOrder, usecaseSession, logger, metrics)

	mux.HandleFunc("/api/v1/order/pay/url", deliveryPayment.OrderGetPayUrl).Methods("GET").Name("get-pay-url")
}
