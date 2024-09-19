package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"
)

func Setup(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, m *metrics.Metrics, reg *prometheus.Registry, logger *zap.Logger) http.Handler {
	logger.Info("begin of handlers definition")
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	AddAuthRouter(mux, cluster, clients, logger, m)
	AddUserRouter(mux, cluster, clients, logger, m)
	AddRestRouter(mux, cluster, clients, logger, m)
	AddOrderRouter(mux, cluster, clients, logger, m)
	AddQuizRouter(mux, cluster, clients, logger, m)
	AddPaymentRouter(mux, cluster, clients, logger, m)

	handler := AddMiddleware(mux, cluster, clients, logger, m)
	logger.Info("end of handlers definition")
	return handler
}
