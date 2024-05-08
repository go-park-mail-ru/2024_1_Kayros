package route

import (
	"database/sql"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/payment"
	rFood "2024_1_kayros/internal/repository/food"
	rOrder "2024_1_kayros/internal/repository/order"
	rRest "2024_1_kayros/internal/repository/restaurants"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/usecase/order"
	ucSession "2024_1_kayros/internal/usecase/session"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

func AddPaymentRouter(db *sql.DB, clientRedisSession *redis.Client, clientRedisCsrf *redis.Client, mux *mux.Router, logger *zap.Logger, cfg *config.Payment) {
	repoFood := rFood.NewRepoLayer(db)
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(clientRedisSession)
	repoCsrf := rSession.NewRepoLayer(clientRedisCsrf)
	repoOrder := rOrder.NewRepoLayer(db)
	repoRest := rRest.NewRepoLayer(db)

	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrf, logger)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseOrder := order.NewUsecaseLayer(repoOrder, repoFood, repoUser, repoRest)
	deliveryPayment := payment.NewPaymentDelivery(logger, usecaseOrder, usecaseCsrf, usecaseSession, cfg)

	mux.HandleFunc("/order/pay/url", deliveryPayment.OrderGetPayUrl).Methods("GET").Name("get-pay-url")
}
