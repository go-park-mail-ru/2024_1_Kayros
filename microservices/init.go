package microservices

import (
	"fmt"

	cfg "2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	grpcClientMiddleware "2024_1_kayros/internal/middleware/grpc/client"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	RestConn    *grpc.ClientConn
	CommentConn *grpc.ClientConn
	AuthConn    *grpc.ClientConn
	UserConn    *grpc.ClientConn
	SessionConn *grpc.ClientConn
}

func Init(logger *zap.Logger, m *metrics.Metrics) *Clients {
	projectCfg := cfg.Config
	middleware := grpcClientMiddleware.NewGrpcClientUnaryMiddlewares(m)
	//restaurant microservice
	restConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", projectCfg.RestGrpcServer.Host, projectCfg.RestGrpcServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(middleware.AccessMiddleware))
	if err != nil {
		logger.Fatal(fmt.Sprintf("the microservice 'restaurant' is not available: %v", err))
	}

	//comment microservice
	commentConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", projectCfg.CommentGrpcServer.Host, projectCfg.CommentGrpcServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(middleware.AccessMiddleware))
	if err != nil {
		logger.Fatal(fmt.Sprintf("the microservice 'comment' is not available: %v", err))
	}

	//auth microservice
	authConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", projectCfg.AuthGrpcServer.Host, projectCfg.AuthGrpcServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(middleware.AccessMiddleware))
	if err != nil {
		logger.Fatal(fmt.Sprintf("the microservice 'authorization' is not available: %v", err))
	}

	// user microservice
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", projectCfg.UserGrpcServer.Host, projectCfg.UserGrpcServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(middleware.AccessMiddleware))
	if err != nil {
		logger.Fatal(fmt.Sprintf("the microservice 'user' is not available: %v", err))
	}

	// session microservice
	sessionConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", projectCfg.SessionGrpcServer.Host, projectCfg.SessionGrpcServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(middleware.AccessMiddleware))
	if err != nil {
		logger.Fatal(fmt.Sprintf("the microservice 'session' is not available: %v", err))
	}

	return &Clients{
		RestConn:    restConn,
		CommentConn: commentConn,
		AuthConn:    authConn,
		UserConn:    userConn,
		SessionConn: sessionConn,
	}
}
