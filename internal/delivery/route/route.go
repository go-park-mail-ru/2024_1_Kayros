package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	metrics "2024_1_kayros"

	"2024_1_kayros/config"
)

func Setup(cfg *config.Project, db *sql.DB, minio *minio.Client, mux *mux.Router, logger *zap.Logger, 
	restConn, commentConn, authConn, userConn, sessionConn *grpc.ClientConn, m *metrics.Metrics) http.Handler {
	logger.Info("The begin of handlers definition")
	mux = mux.PathPrefix("/api/v1").Subrouter()
	mux.StrictSlash(true)
	mux.PathPrefix("/metrics").Handler(promhttp.Handler())

	AddAuthRouter(cfg, db, authConn, sessionConn, mux, logger)
	AddUserRouter(db, cfg, userConn, sessionConn, mux, logger)
	AddRestRouter(db, mux, logger, restConn, commentConn)
	AddOrderRouter(db, mux, logger)
	AddQuizRouter(db, sessionConn, userConn, minio, mux, logger, cfg)
	AddPaymentRouter(db, sessionConn, mux, logger, cfg)

	handler := AddMiddleware(cfg, db, sessionConn, userConn, mux, logger, m)
	logger.Info("The end of handlers definition")
	return handler
}
