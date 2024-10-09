package client

import (
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/utils/constants"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcClientUnaryMiddlewares struct {
	metrics *metrics.Metrics
}

func NewGrpcClientUnaryMiddlewares(metrics *metrics.Metrics) *GrpcClientUnaryMiddlewares {
	return &GrpcClientUnaryMiddlewares{
		metrics: metrics,
	}
}

func (mdlwr *GrpcClientUnaryMiddlewares) AccessMiddleware(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	requestId := ""
	value := ctx.Value(constants.RequestId)
	if value != nil {
		requestId = value.(string)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, constants.RequestId, requestId)

	return invoker(ctx, method, req, reply, cc, opts...)
}
