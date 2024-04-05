package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Setup(db *sql.DB, redis *redis.Client, minio *minio.Client, mux *mux.Router, logger *zap.Logger) {
	logger.Info("Начало определения хендлеров")
	mux.PathPrefix("/api/v1")
	mux.StrictSlash(true)

	AddAuthRouter(db, redis, minio, mux, logger)
	AddUserRouter(db, minio, mux, logger)
	AddRestRouter(db, mux, logger)
	AddOrderRouter(db, minio, mux, logger)

	AddMiddleware(db, redis, minio, mux, logger)
	logger.Info("Конец определения хендлеров")
}
