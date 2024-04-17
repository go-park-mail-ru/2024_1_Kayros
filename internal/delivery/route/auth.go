package route

import (
	"database/sql"

	"2024_1_kayros/config"
	dAuth "2024_1_kayros/internal/delivery/auth"
	rMinio "2024_1_kayros/internal/repository/minios3"
	rSession "2024_1_kayros/internal/repository/session"
	rUser "2024_1_kayros/internal/repository/user"
	ucAuth "2024_1_kayros/internal/usecase/auth"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AddAuthRouter(cfg *config.Project, db *sql.DB, clientRedisSession *redis.Client, clientRedisCsrf *redis.Client, clientMinio *minio.Client, mux *mux.Router, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(clientRedisSession)
	repoCsrf := rSession.NewRepoLayer(clientRedisCsrf)
	repoMinio := rMinio.NewRepoLayer(clientMinio)

	usecaseUser := ucUser.NewUsecaseLayer(repoUser, repoMinio)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrf, logger)
	usecaseAuth := ucAuth.NewUsecaseLayer(repoUser)

	deliveryAuth := dAuth.NewDeliveryLayer(cfg, usecaseSession, usecaseUser, usecaseCsrf, usecaseAuth, logger)

	mux.HandleFunc("/signin", deliveryAuth.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/signup", deliveryAuth.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/signout", deliveryAuth.SignOut).Methods("POST").Name("signout")
}
