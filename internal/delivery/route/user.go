package route

import (
	"database/sql"

	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/gen/go/user"
	dUser "2024_1_kayros/internal/delivery/user"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func AddUserRouter(db *sql.DB, cfg *config.Project, userConn, sessionConn *grpc.ClientConn, mux *mux.Router, logger *zap.Logger) {
	// init user grpc client
	grpcUserClient := user.NewUserManagerClient(userConn)
	usecaseUser := ucUser.NewUsecaseLayer(grpcUserClient)
	// init session grpc client
	grpcSessionClient := session.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcSessionClient)

	deliveryUser := dUser.NewDeliveryLayer(cfg, usecaseSession, usecaseUser, logger)

	mux.HandleFunc("/api/v1/user", deliveryUser.UserData).Methods("GET").Name("user_data")
	mux.HandleFunc("/api/v1/user", deliveryUser.UpdateInfo).Methods("PUT").Name("update_user")
	mux.HandleFunc("/api/v1/user/address", deliveryUser.UpdateAddress).Methods("PUT").Name("update_user_address")
	mux.HandleFunc("/api/v1/user/address", deliveryUser.UserAddress).Methods("GET").Name("user_address")
	mux.HandleFunc("/api/v1/user/new_password", deliveryUser.UpdatePassword).Methods("PUT").Name("update_user_password")
}
