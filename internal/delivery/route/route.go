package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func Setup(db *sql.DB, redis *redis.Client, mux *mux.Router) {
	mux.PathPrefix("/api/v1")
	mux.StrictSlash(true)

	AddAuthRouter(db, redis, mux)
	AddRestRouter(db, mux)
	AddUserRouter(db, mux)

	AddMiddleware(db, redis, mux)
}
