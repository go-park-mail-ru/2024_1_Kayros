package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Setup(db *sql.DB, redis *redis.Client, minio *minio.Client, mux *mux.Router, logger *zap.Logger) http.Handler {
	logger.Info("Начало определения хендлеров")
	mux = mux.PathPrefix("/api/v1").Subrouter()
	mux.StrictSlash(true)

	AddAuthRouter(db, redis, minio, mux, logger) // протестировано
	AddUserRouter(db, minio, mux, logger)        // протестировано
	AddRestRouter(db, mux, logger)
	AddOrderRouter(db, minio, mux, logger)

	handler := AddMiddleware(db, redis, minio, mux, logger)
	logger.Info("Конец определения хендлеров")
	return handler
}
