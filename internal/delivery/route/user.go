package route

import (
	"database/sql"

	"2024_1_kayros/config"
	dUser "2024_1_kayros/internal/delivery/user"
	rMinio "2024_1_kayros/internal/repository/minios3"
	rSession "2024_1_kayros/internal/repository/session"
	ucSession "2024_1_kayros/internal/usecase/session"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	rUser "2024_1_kayros/internal/repository/user"

	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
)

<<<<<<< HEAD
func AddUserRouter(db *sql.DB, cfg *config.Project, minio *minio.Client, clientRedisSession *redis.Client, clientRedisCsrf *redis.Client, clientRedisUnauthTokens *redis.Client, mux *mux.Router, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(db, minio, logger)
	repoSession := rSession.NewRepoLayer(clientRedisSession, logger)
	repoCsrf := rSession.NewRepoLayer(clientRedisCsrf, logger)
	repoUnauthAddress := rSession.NewRepoLayer(clientRedisUnauthTokens, logger)
=======
func AddUserRouter(db *sql.DB, cfg *config.Project, minio *minio.Client, clientRedisSession *redis.Client, clientRedisCsrf *redis.Client, mux *mux.Router, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(db)
	repoSession := rSession.NewRepoLayer(clientRedisSession)
	repoCsrf := rSession.NewRepoLayer(clientRedisCsrf)
	repoMinio := rMinio.NewRepoLayer(minio)
>>>>>>> fix_csrf_test

	usecaseUser := ucUser.NewUsecaseLayer(repoUser, repoMinio)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrf, logger)
	usecaseUnauthAddress := ucSession.NewUsecaseLayer(repoUnauthAddress, logger)
	deliveryUser := dUser.NewDeliveryLayer(cfg, usecaseSession, usecaseUser, usecaseCsrf, usecaseUnauthAddress, logger)

	mux.HandleFunc("/user", deliveryUser.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("/user", deliveryUser.UpdateInfo).Methods("PUT").Name("update_user")
	mux.HandleFunc("/user/address", deliveryUser.UpdateAddress).Methods("PUT").Name("update_user_address")
	mux.HandleFunc("/user/new_password", deliveryUser.UpdatePassword).Methods("PUT").Name("update_user_password")
}
