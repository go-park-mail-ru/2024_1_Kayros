package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/internal/delivery/metrics"

	"2024_1_kayros/config"
)


func Setup(cfg *config.Project, db *sql.DB, minio *minio.Client, statements map[string]map[string]*sql.Stmt, mux *mux.Router, logger *zap.Logger,
	restConn, commentConn, authConn, userConn, sessionConn *grpc.ClientConn, m *metrics.Metrics, reg *prometheus.Registry) http.Handler {
	logger.Info("The begin of handlers definition")
	mux.StrictSlash(true)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)
	
	AddAuthRouter(cfg, authConn, userConn, sessionConn, mux, logger, m)
	AddUserRouter(cfg, userConn, sessionConn, mux, logger, m)
	AddRestRouter(db, statements, mux, logger, restConn, userConn, commentConn, m)
	AddOrderRouter(db, statements, mux, userConn, restConn, logger, m)
	AddQuizRouter(db, statements["statistic"], sessionConn, userConn, minio, mux, logger, cfg, m)
	AddPaymentRouter(db, statements, sessionConn, userConn, restConn, mux, logger, cfg, m)

	handler := AddMiddleware(cfg, db, sessionConn, userConn, mux, logger, m)
	logger.Info("The end of handlers definition")
	return handler
}
