package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"2024_1_kayros/config"
)

func Setup(cfg *config.Project, db *sql.DB, redisSession *redis.Client, redisCsrf *redis.Client, minio *minio.Client, mux *mux.Router, logger *zap.Logger, restConn, commentConn *grpc.ClientConn) http.Handler {
	logger.Info("The begin of handlers definition")
	mux = mux.PathPrefix("/api/v1").Subrouter()
	mux.StrictSlash(true)

	AddAuthRouter(cfg, db, redisSession, redisCsrf, minio, mux, logger)
	AddUserRouter(db, cfg, minio, redisSession, redisCsrf, mux, logger)
	AddRestRouter(db, mux, logger, restConn)
	AddOrderRouter(db, mux, logger)
	AddQuizRouter(db, redisSession, redisCsrf, minio, mux, logger)
	AddCommentRouter(db, mux, logger, commentConn)

	handler := AddMiddleware(cfg, db, redisSession, redisCsrf, minio, mux, logger)
	logger.Info("The end of handlers definition")
	return handler
}
