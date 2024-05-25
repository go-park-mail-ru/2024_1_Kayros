package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/config"
	protosession "2024_1_kayros/gen/go/session"
	protouser "2024_1_kayros/gen/go/user"
	http_middleware "2024_1_kayros/internal/middleware/http"
	ucSession "2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
)

func AddMiddleware(cfg *config.Project, db *sql.DB, sessionConn, userConn *grpc.ClientConn, mux *mux.Router, logger *zap.Logger, m *metrics.Metrics) http.Handler {
	// init user microservice client
	grpcUserClient := protouser.NewUserManagerClient(userConn)
	usecaseUser := user.NewUsecaseLayer(grpcUserClient, m)
	// init session microservice client
	grpcSessionClient := protosession.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient, m)

	handler := http_middleware.SessionAuthentication(mux, usecaseUser, usecaseSession, logger, &cfg.Redis)
	handler = http_middleware.Csrf(handler, usecaseSession, cfg, logger, m)
	handler = http_middleware.Cors(handler, logger)
	handler = http_middleware.Access(handler, logger)
	handler = http_middleware.Metrics(handler, m)

	return handler
}
