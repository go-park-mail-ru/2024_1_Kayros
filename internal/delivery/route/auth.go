package route

import (
	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/gen/go/user"
	dAuth "2024_1_kayros/internal/delivery/auth"
	"2024_1_kayros/internal/delivery/metrics"
	ucAuth "2024_1_kayros/internal/usecase/auth"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func AddAuthRouter(cfg *config.Project, authConn *grpc.ClientConn, userConn *grpc.ClientConn, sessionConn *grpc.ClientConn,
	mux *mux.Router, logger *zap.Logger, metrics *metrics.Metrics) {
	// microservice authorization
	grpcClientAuth := auth.NewAuthManagerClient(authConn)
	usecaseAuth := ucAuth.NewUsecaseLayer(grpcClientAuth, metrics)
	// microservice session
	grpcClientSession := session.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcClientSession, metrics)
	// microservice user
	grpcClientUser := user.NewUserManagerClient(userConn)
	usecaseUser := ucUser.NewUsecaseLayer(grpcClientUser, metrics)

	deliveryAuth := dAuth.NewDeliveryLayer(cfg, usecaseSession, usecaseAuth, usecaseUser, logger, metrics)

	mux.HandleFunc("/api/v1/vk", deliveryAuth.AuthVk).Methods("POST").Name("vk-auth")
	mux.HandleFunc("/api/v1/signin", deliveryAuth.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/api/v1/signup", deliveryAuth.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/api/v1/signout", deliveryAuth.SignOut).Methods("POST").Name("signout")
}

