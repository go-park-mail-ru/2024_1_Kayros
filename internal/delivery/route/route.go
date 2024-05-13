package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"

	"2024_1_kayros/internal/delivery/metrics"

	"2024_1_kayros/config"
)


func Setup(cfg *config.Project, db *sql.DB, minio *minio.Client, mux *mux.Router, logger *zap.Logger,
	restConn, commentConn, authConn, userConn, sessionConn *grpc.ClientConn, m *metrics.Metrics, reg *prometheus.Registry) http.Handler {
	logger.Info("The begin of handlers definition")
	mux.StrictSlash(true)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)
	
	AddAuthRouter(cfg, db, authConn, sessionConn, mux, logger)
	AddUserRouter(db, cfg, userConn, sessionConn, mux, logger)
	AddRestRouter(db, mux, logger, restConn, userConn, commentConn)
	AddOrderRouter(db, mux, userConn, restConn, logger)
	AddQuizRouter(db, sessionConn, userConn, minio, mux, logger, cfg)
	AddPaymentRouter(db, sessionConn, userConn, restConn, mux, logger, cfg)

	handler := AddMiddleware(cfg, db, sessionConn, userConn, mux, logger, m)
	logger.Info("The end of handlers definition")
	return handler
}
