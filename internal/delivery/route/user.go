package route

import (
	"database/sql"

	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/gen/go/session"
	dUser "2024_1_kayros/internal/delivery/user"
	ucUser "2024_1_kayros/internal/usecase/user"
	ucSession "2024_1_kayros/internal/usecase/session"

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

	mux.HandleFunc("/user", deliveryUser.UserData).Methods("GET").Name("user_data")
	mux.HandleFunc("/user", deliveryUser.UpdateInfo).Methods("PUT").Name("update_user")
	mux.HandleFunc("/user/address", deliveryUser.UpdateAddress).Methods("PUT").Name("update_user_address")
	mux.HandleFunc("/user/address", deliveryUser.UserAddress).Methods("GET").Name("user_address")
	mux.HandleFunc("/user/new_password", deliveryUser.UpdatePassword).Methods("PUT").Name("update_user_password")
}
