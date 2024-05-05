package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	metrics "2024_1_kayros"
	"2024_1_kayros/config"
	"2024_1_kayros/internal/middleware"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
)

func AddMiddleware(cfg *config.Project, db *sql.DB, redisClientSession *redis.Client, redisClientCsrf *redis.Client, mux *mux.Router, logger *zap.Logger, m *metrics.Metrics) http.Handler {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(redisClientSession)
	repoCsrfTokens := rSession.NewRepoLayer(redisClientCsrf)

	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrfTokens, logger)

	handler := middleware.SessionAuthentication(mux, repoUser, usecaseSession, logger)
	handler = middleware.Csrf(handler, usecaseCsrf, usecaseSession, cfg, logger)
	handler = middleware.Cors(handler, logger)
	handler = middleware.Access(handler, logger)
	handler = middleware.Metrics(handler, m)

	return handler
}
