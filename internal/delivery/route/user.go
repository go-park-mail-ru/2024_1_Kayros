package route

import (
	"database/sql"

	dUser "2024_1_kayros/internal/delivery/user"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	rUser "2024_1_kayros/internal/repository/user"

	ucUser "2024_1_kayros/internal/usecase/user"
	"github.com/gorilla/mux"
)

func AddUserRouter(db *sql.DB, minio *minio.Client, mux *mux.Router, logger *zap.Logger) {
	repoUser := rUser.NewRepoLayer(db, minio, logger)
	usecaseUser := ucUser.NewUsecaseLayer(repoUser, logger)
	deliveryUser := dUser.NewDeliveryLayer(usecaseUser, logger)

	mux.HandleFunc("user", deliveryUser.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("user", deliveryUser.UploadImage).Methods("PUT").Name("user_image")
}
