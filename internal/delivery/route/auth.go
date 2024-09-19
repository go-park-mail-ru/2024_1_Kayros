package route

import (
	genAuth "2024_1_kayros/gen/go/auth"
	genSession "2024_1_kayros/gen/go/session"
	genUser "2024_1_kayros/gen/go/user"
	dAuth "2024_1_kayros/internal/delivery/auth"
	"2024_1_kayros/internal/delivery/metrics"
	ucAuth "2024_1_kayros/internal/usecase/auth"
	ucSession "2024_1_kayros/internal/usecase/session"
	ucUser "2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/microservices"
	"2024_1_kayros/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func AddAuthRouter(mux *mux.Router, cluster *services.Cluster, clients *microservices.Clients, logger *zap.Logger, metrics *metrics.Metrics) {
	// microservice authorization
	grpcClientAuth := genAuth.NewAuthManagerClient(clients.AuthConn)
	usecaseAuth := ucAuth.NewUsecaseLayer(grpcClientAuth, metrics)
	// microservice session
	grpcClientSession := genSession.NewSessionManagerClient(clients.SessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcClientSession, metrics)
	// microservice user
	grpcClientUser := genUser.NewUserManagerClient(clients.UserConn)
	usecaseUser := ucUser.NewUsecaseLayer(grpcClientUser, metrics)

	deliveryAuth := dAuth.NewDeliveryLayer(usecaseSession, usecaseAuth, usecaseUser, logger, metrics)

	mux.HandleFunc("/api/v1/vk", deliveryAuth.AuthVk).Methods("POST").Name("vk-auth")
	mux.HandleFunc("/api/v1/signin", deliveryAuth.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/api/v1/signup", deliveryAuth.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/api/v1/signout", deliveryAuth.SignOut).Methods("POST").Name("signout")
}
