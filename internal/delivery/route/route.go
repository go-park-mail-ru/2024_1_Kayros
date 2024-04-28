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

<<<<<<< HEAD
func Setup(cfg *config.Project, db *sql.DB, redisSession *redis.Client, redisCsrf *redis.Client, redisUnauthTokens *redis.Client, minio *minio.Client, mux *mux.Router, logger *zap.Logger) http.Handler {
	logger.Info("Начало определения хендлеров")
=======
func Setup(cfg *config.Project, db *sql.DB, redisSession *redis.Client, redisCsrf *redis.Client, minio *minio.Client, mux *mux.Router, logger *zap.Logger) http.Handler {
	logger.Info("The begin of handlers definition")
>>>>>>> fix_csrf_test
	mux = mux.PathPrefix("/api/v1").Subrouter()
	mux.StrictSlash(true)

	AddAuthRouter(cfg, db, redisSession, redisCsrf, redisUnauthTokens, minio, mux, logger) // протестировано
	AddUserRouter(db, cfg, minio, redisSession, redisCsrf, redisUnauthTokens, mux, logger) // протестировано
	AddRestRouter(db, mux, logger)
	AddOrderRouter(db, minio, mux, logger)
	AddQuizRouter(db, minio, mux, logger)

	handler := AddMiddleware(cfg, db, redisSession, redisCsrf, minio, mux, logger)
	logger.Info("The end of handlers definition")
	return handler
}
