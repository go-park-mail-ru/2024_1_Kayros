package route

import (
	"database/sql"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/middleware"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AddMiddleware(cfg *config.Project, db *sql.DB, redisClientSession *redis.Client, redisClientCsrf *redis.Client, minioClient *minio.Client, mux *mux.Router, logger *zap.Logger) http.Handler {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(redisClientSession)
	repoCsrfTokens := rSession.NewRepoLayer(redisClientCsrf)

	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrfTokens, logger)

	handler := middleware.SessionAuthentication(mux, repoUser, usecaseSession, logger)
	handler = middleware.Csrf(handler, usecaseCsrf, usecaseSession, cfg, logger)
	handler = middleware.Cors(handler, logger)
	handler = middleware.Access(handler, logger)
	return handler
}
