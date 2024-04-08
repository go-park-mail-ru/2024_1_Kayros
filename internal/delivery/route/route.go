package route

import (
	"database/sql"
	"net/http"

	"2024_1_kayros/config"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Setup(cfg *config.Project, db *sql.DB, redisSession *redis.Client, redisCsrf *redis.Client, minio *minio.Client, mux *mux.Router, logger *zap.Logger) http.Handler {
	logger.Info("Начало определения хендлеров")
	mux = mux.PathPrefix("/api/v1").Subrouter()
	mux.StrictSlash(true)

	AddAuthRouter(cfg, db, redisSession, redisCsrf, minio, mux, logger) // протестировано
	AddUserRouter(db, minio, mux, logger)                               // протестировано
	AddRestRouter(db, mux, logger)
	AddOrderRouter(db, minio, mux, logger)

	handler := AddMiddleware(cfg, db, redisSession, redisCsrf, minio, mux, logger)
	logger.Info("Конец определения хендлеров")
	return handler
}
