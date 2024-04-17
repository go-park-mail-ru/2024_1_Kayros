package route

import (
	"database/sql"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/middleware"
	"2024_1_kayros/internal/repository/minios3"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AddMiddleware(cfg *config.Project, db *sql.DB, redisClientSession *redis.Client, redisClientCsrf *redis.Client, minioClient *minio.Client, mux *mux.Router, logger *zap.Logger) http.Handler {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(redisClientSession)
	repoCsrfTokens := rSession.NewRepoLayer(redisClientCsrf)
	repoMinio := minios3.NewRepoLayer(minioClient)

	usecaseUser := ucUser.NewUsecaseLayer(repoUser, repoMinio)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrfTokens, logger)

	// цепочка middlewares
	handler := middleware.CsrfMiddleware(mux, usecaseCsrf, cfg, logger)
	handler = middleware.SessionAuthenticationMiddleware(handler, usecaseUser, usecaseSession, logger)
	handler = middleware.CorsMiddleware(handler, logger)
	handler = middleware.AccessMiddleware(handler, logger)
	return handler
}
