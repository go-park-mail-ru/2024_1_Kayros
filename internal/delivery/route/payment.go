package route

import (
	"database/sql"

	"2024_1_kayros/config"
	restproto "2024_1_kayros/gen/go/rest"
	sessionproto "2024_1_kayros/gen/go/session"
	userproto "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/delivery/payment"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/usecase/order"
	ucSession "2024_1_kayros/internal/usecase/session"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

func AddPaymentRouter(db *sql.DB, statements map[string]map[string]*sql.Stmt, sessionConn, userConn, restConn *grpc.ClientConn, mux *mux.Router, logger *zap.Logger, cfg *config.Project, metrics *metrics.Metrics) {
	repoFood := rFood.NewLayer(db, statements["food"])
	repoOrder := rOrder.NewRepoLayer(db, metrics, statements["order"])

	// init session grpc client
	grpcSessionClient := sessionproto.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient, metrics)
	// init user grpc client
	grpcUserClient := userproto.NewUserManagerClient(userConn)
	// init rest grpc client
	grpcRestClient := restproto.NewRestWorkerClient(restConn)

	usecaseOrder := order.NewUsecaseLayer(repoOrder, repoFood, grpcUserClient, grpcRestClient, metrics)
	deliveryPayment := payment.NewPaymentDelivery(logger, usecaseOrder, usecaseSession, cfg, metrics)

	mux.HandleFunc("/api/v1/order/pay/url", deliveryPayment.OrderGetPayUrl).Methods("GET").Name("get-pay-url")
}
