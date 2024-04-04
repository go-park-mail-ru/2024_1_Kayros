package route

import (
	"database/sql"

	"2024_1_kayros/internal/middleware"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AddMiddleware(db *sql.DB, redisClient *redis.Client, minioClient *minio.Client, mux *mux.Router, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(db, minioClient, logger)
	repoSession := rSession.NewRepoLayer(redisClient, logger)

	usecaseUser := ucUser.NewUsecaseLayer(repoUser, logger)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)

	// цепочка middlewares
	handler := middleware.SessionAuthenticationMiddleware(mux, usecaseUser, usecaseSession, logger)
	handler = middleware.CorsMiddleware(handler, logger)
}
