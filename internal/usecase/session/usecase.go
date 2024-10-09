package session

import (
	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=./usecase.go -destination=./usecase_mock.go -package=session
type Usecase interface {
	GetValue(ctx context.Context, key alias.SessionKey, databaseNum int32) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, databaseNum int32) error
	DeleteKey(ctx context.Context, key alias.SessionKey, databaseNum int32) error
}

type UsecaseLayer struct {
	client  session.SessionManagerClient
	metrics *metrics.Metrics
}

func NewUsecaseLayer(clientProps session.SessionManagerClient, metrics *metrics.Metrics) Usecase {
	return &UsecaseLayer{
		client:  clientProps,
		metrics: metrics,
	}
}

func (uc *UsecaseLayer) GetValue(ctx context.Context, key alias.SessionKey, databaseNum int32) (alias.SessionValue, error) {
	data := &session.GetSessionData{
		Key:      string(key),
		Database: databaseNum,
	}
	timeNow := time.Now()
	value, err := uc.client.GetSession(ctx, data)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.SessionMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.SessionMicroservice, grpcStatus.String()).Inc()
		}
		if grpcerr.Is(err, codes.NotFound, myerrors.RedisNoData) {
			return "", myerrors.RedisNoData
		}
		return "", err
	}
	return alias.SessionValue(value.GetData()), err
}

func (uc *UsecaseLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, databaseNum int32) error {
	data := &session.SetSessionData{
		Key:      string(key),
		Value:    string(value),
		Database: databaseNum,
	}
	timeNow := time.Now()
	_, err := uc.client.SetSession(ctx, data)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.SessionMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.SessionMicroservice, grpcStatus.String()).Inc()
		}
		if grpcerr.Is(err, codes.Internal, myerrors.RedisNoData) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

func (uc *UsecaseLayer) DeleteKey(ctx context.Context, key alias.SessionKey, databaseNum int32) error {
	data := &session.DeleteSessionData{
		Key:      string(key),
		Database: databaseNum,
	}
	timeNow := time.Now()
	_, err := uc.client.DeleteSession(ctx, data)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.SessionMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.SessionMicroservice, grpcStatus.String()).Inc()
		}
		if grpcerr.Is(err, codes.Internal, myerrors.RedisNoData) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}
