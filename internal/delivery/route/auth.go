package route

import (
	"database/sql"

	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/session"
	dAuth "2024_1_kayros/internal/delivery/auth"
	"2024_1_kayros/internal/delivery/metrics"
	ucAuth "2024_1_kayros/internal/usecase/auth"
	ucSession "2024_1_kayros/internal/usecase/session"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func AddAuthRouter(cfg *config.Project, db *sql.DB, authConn *grpc.ClientConn, sessionConn *grpc.ClientConn,
	mux *mux.Router, logger *zap.Logger, metrics *metrics.Metrics) {
	// microservice authorization
	grpcClientAuth := auth.NewAuthManagerClient(authConn)
	usecaseAuth := ucAuth.NewUsecaseLayer(grpcClientAuth, metrics)
	// microservice session
	grpcClientSession := session.NewSessionManagerClient(sessionConn)
	usecaseSession := ucSession.NewUsecaseLayer(grpcClientSession, metrics)

	deliveryAuth := dAuth.NewDeliveryLayer(cfg, usecaseSession, usecaseAuth, logger, metrics)

	mux.HandleFunc("/api/v1/signin", deliveryAuth.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/api/v1/signup", deliveryAuth.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/api/v1/signout", deliveryAuth.SignOut).Methods("POST").Name("signout")
}

