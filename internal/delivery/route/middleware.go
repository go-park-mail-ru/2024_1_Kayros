package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	metrics "2024_1_kayros"
	"2024_1_kayros/config"
	protouser "2024_1_kayros/gen/go/user"
	protosession "2024_1_kayros/gen/go/session"
	"2024_1_kayros/internal/middleware"
	ucSession "2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
)

func AddMiddleware(cfg *config.Project, db *sql.DB, sessionConn, userConn *grpc.ClientConn, mux *mux.Router, logger *zap.Logger, m *metrics.Metrics) http.Handler {
	// init user microservice client
	grpcUserClient := protouser.NewUserManagerClient(userConn)
	usecaseUser := user.NewUsecaseLayer(grpcUserClient)
	// init session microservice client
	grpcSessionClient := protosession.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient)

	handler := middleware.SessionAuthentication(mux, usecaseUser, usecaseSession, logger, &cfg.Redis)
	handler = middleware.Csrf(handler, usecaseSession, cfg, logger)
	handler = middleware.Cors(handler, logger)
	handler = middleware.Access(handler, logger)
	handler = middleware.Metrics(handler, m)

	return handler
}
