package rest

import (
	"context"
	"strconv"
	"time"

	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"

	"google.golang.org/grpc/status"
)

type Usecase interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
	GetByFilter(ctx context.Context, id alias.CategoryId) ([]*entity.Restaurant, error)
	GetCategoryList(ctx context.Context) ([]*entity.Category, error)
	GetRecomendation(ctx context.Context, userId uint64) ([]*entity.Restaurant, error)
}
type UsecaseLayer struct {
	grpcClient rest.RestWorkerClient
	metrics    *metrics.Metrics
}

func NewUsecaseLayer(restClient rest.RestWorkerClient, m *metrics.Metrics) *UsecaseLayer {
	return &UsecaseLayer{
		grpcClient: restClient,
		metrics:    m,
	}
}

func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	timeNow := time.Now()
	rests, err := uc.grpcClient.GetAll(ctx, nil)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.RestMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.RestMicroservice, grpcStatus.String()).Inc()
		}
		return nil, err
	}
	return FromGrpcStructToRestaurantArray(rests), nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	timeNow := time.Now()
	r, err := uc.grpcClient.GetById(ctx, &rest.RestId{Id: uint64(restId)})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.RestMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	uc.metrics.PopularRestaurant.WithLabelValues(strconv.Itoa(int(restId))).Inc()
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.RestMicroservice, grpcStatus.String()).Inc()
		}
		return nil, err
	}
	return FromGrpcStructToRestaurant(r), nil
}

func (uc *UsecaseLayer) GetByFilter(ctx context.Context, id alias.CategoryId) ([]*entity.Restaurant, error) {
	timeNow := time.Now()
	rests, err := uc.grpcClient.GetByFilter(ctx, &rest.Id{Id: uint64(id)})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.RestMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	uc.metrics.PopularRestaurant.WithLabelValues(strconv.Itoa(int(id))).Inc()
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.RestMicroservice, grpcStatus.String()).Inc()
		}
		return nil, err
	}
	return FromGrpcStructToRestaurantArray(rests), nil
}

func (uc *UsecaseLayer) GetCategoryList(ctx context.Context) ([]*entity.Category, error) {
	timeNow := time.Now()
	cats, err := uc.grpcClient.GetCategoryList(ctx, nil)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.RestMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.RestMicroservice, grpcStatus.String()).Inc()
		}
		return nil, err
	}
	return FromGrpcStructToCategoryArray(cats), nil
}

func (uc *UsecaseLayer) GetRecomendation(ctx context.Context, userId uint64) ([]*entity.Restaurant, error) {
	rests, err := uc.grpcClient.GetRecomendation(ctx, &rest.UserAndLimit{UserId: userId, Limit: 5})
	if err != nil {
		return nil, err
	}
	return FromGrpcStructToRestaurantArray(rests), nil
}
