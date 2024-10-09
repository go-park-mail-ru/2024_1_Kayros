package server

import (
	"2024_1_kayros/internal/utils/constants"
	metrics "2024_1_kayros/microservices/metrics"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	responseGrpcStatus = "response_grpc_status"
	latencyHumanMs     = "latency_human_ms"
	databaseDurationMs = "database_duration_ms"
	repoMethod         = "repo_method"
)

type MiddlewareChain struct {
	logger  *zap.Logger
	metrics *metrics.MicroserviceMetrics
}

func NewMiddlewareChain(logger *zap.Logger, metrics *metrics.MicroserviceMetrics) MiddlewareChain {
	return MiddlewareChain{
		logger:  logger,
		metrics: metrics,
	}
}

func (mdlwr *MiddlewareChain) AccessMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestId := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md[constants.RequestId]) > 0 {
		requestId = md[constants.RequestId][0]
	}
	mdlwr.logger.Info(fmt.Sprintf("init request %s", requestId))

	timeNow := time.Now()
	resp, err := handler(ctx, req)
	timeEnd := time.Since(timeNow).String()

	//receiving code status
	grpcCode := codes.OK
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			mdlwr.logger.Error("Code status is unavaliable. Status code --> " + grpcStatus.String())
		}
	}

	mdlwr.logger.Info(fmt.Sprintf("request done %s", requestId),
		zap.String(latencyHumanMs, timeEnd),
		zap.String(responseGrpcStatus, grpcCode.String()))
	return resp, err
}

func (mdlwr *MiddlewareChain) MetricsMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	timeNow := time.Now()
	resp, err := handler(ctx, req)
	timeEnd := time.Since(timeNow)

	//receiving code status
	grpcCode := codes.OK
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			mdlwr.logger.Error("Code status is unavaliable. Status code --> " + grpcStatus.String())
		}
	}
	//increment number of requests
	mdlwr.metrics.TotalNumberOfRequests.Inc()
	//add status and time of request
	mdlwr.metrics.RequestTime.WithLabelValues(grpcCode.String()).Observe(float64(timeEnd.Milliseconds()))
	return resp, err
}
