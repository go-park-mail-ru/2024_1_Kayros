package rest

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

type Usecase interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
}
type UsecaseLayer struct {
	grpcClient rest.RestWorkerClient
}

func NewUsecaseLayer(restClient rest.RestWorkerClient) *UsecaseLayer {
	return &UsecaseLayer{
		grpcClient: restClient,
	}
}

func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	rests, err := uc.grpcClient.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	return entity.FromGrpcStructToRestaurantArray(rests), nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	r, err := uc.grpcClient.GetById(ctx, &rest.RestId{Id: uint64(restId)})
	if err != nil {
		return nil, err
	}
	return entity.FromGrpcStructToRestaurant(r), nil
}
