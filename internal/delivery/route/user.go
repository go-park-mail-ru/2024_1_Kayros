package route

import (
	"database/sql"

	"2024_1_kayros/config"
	dUser "2024_1_kayros/internal/delivery/user"
	rSession "2024_1_kayros/internal/repository/session"
	ucSession "2024_1_kayros/internal/usecase/session"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	rUser "2024_1_kayros/internal/repository/user"

	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
)

func AddUserRouter(db *sql.DB, cfg *config.Project, minio *minio.Client, clientRedisSession *redis.Client, clientRedisCsrf *redis.Client, mux *mux.Router, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(db, minio, logger)
	repoSession := rSession.NewRepoLayer(clientRedisSession, logger)
	repoCsrf := rSession.NewRepoLayer(clientRedisCsrf, logger)

	usecaseUser := ucUser.NewUsecaseLayer(repoUser, logger)
	usecaseSession := ucSession.NewUsecaseLayer(repoSession, logger)
	usecaseCsrf := ucSession.NewUsecaseLayer(repoCsrf, logger)
	deliveryUser := dUser.NewDeliveryLayer(cfg, usecaseSession, usecaseUser, usecaseCsrf, logger)

	mux.HandleFunc("/user", deliveryUser.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("/user", deliveryUser.UpdateInfo).Methods("PUT").Name("update_user")
}
