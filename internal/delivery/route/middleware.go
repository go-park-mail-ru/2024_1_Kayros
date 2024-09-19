package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/config"
	protosession "2024_1_kayros/gen/go/session"
	protouser "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	http_middleware "2024_1_kayros/internal/middleware/http"
	ucSession "2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"
)

func AddMiddleware(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, logger *zap.Logger, m *metrics.Metrics) http.Handler {
	// init user microservice client
	grpcUserClient := protouser.NewUserManagerClient(clients.UserConn)
	usecaseUser := user.NewUsecaseLayer(grpcUserClient, m)
	// init session microservice client
	grpcSessionClient := protosession.NewSessionManagerClient(clients.SessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient, m)

	handler := http_middleware.SessionAuthentication(mux, usecaseUser, usecaseSession, &config.Config.Redis, logger)
	handler = http_middleware.Csrf(handler, usecaseSession, &config.Config, logger, m)
	handler = http_middleware.Cors(handler, logger)
	handler = http_middleware.Access(handler, logger)
	handler = http_middleware.Metrics(handler, m)

	return handler
}
