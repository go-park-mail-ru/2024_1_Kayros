package route

import (
	"database/sql"

	dAuth "2024_1_kayros/internal/delivery/auth"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

// нужно будет добавить интерфейс к БД и редису
func AddAuthRouter(db *sql.DB, client *redis.Client, mux *mux.Router) {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(client)

	usecaseUser := ucUser.NewUsecaseLayer(repoUser)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession)

	deliveryAuth := dAuth.NewDeliveryLayer(usecaseSession, usecaseUser)

	mux.HandleFunc("signin", deliveryAuth.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("signup", deliveryAuth.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("signout", deliveryAuth.SignOut).Methods("POST").Name("signout")
}
