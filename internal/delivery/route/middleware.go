package route

import (
	"database/sql"

	"2024_1_kayros/internal/middleware"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func AddMiddleware(db *sql.DB, client *redis.Client, mux *mux.Router) {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(client)

	usecaseUser := ucUser.NewUsecaseLayer(repoUser)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession)

	// цепочка middlewares
	handler := middleware.SessionAuthenticationMiddleware(mux, usecaseUser, usecaseSession)
	handler = middleware.CorsMiddleware(handler)
}
