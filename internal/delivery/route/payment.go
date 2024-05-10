package route

import (
	"database/sql"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/payment"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	rRest "2024_1_kayros/internal/repository/restaurants"
	"2024_1_kayros/internal/usecase/order"
	ucSession "2024_1_kayros/internal/usecase/session"
	sessionproto "2024_1_kayros/gen/go/session"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

func AddPaymentRouter(db *sql.DB, sessionConn *grpc.ClientConn, mux *mux.Router, logger *zap.Logger, cfg *config.Project) {
	repoFood := rFood.NewRepoLayer(db)
	repoOrder := rOrder.NewRepoLayer(db)

	// init session grpc client
	grpcSessionClient := sessionproto.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient)

	usecaseOrder := order.NewUsecaseLayer(repoOrder, repoFood, repoUser, repoRest)
	deliveryPayment := payment.NewPaymentDelivery(logger, usecaseOrder, usecaseSession, cfg)

	mux.HandleFunc("/order/pay/url", deliveryPayment.OrderGetPayUrl).Methods("GET").Name("get-pay-url")
}
