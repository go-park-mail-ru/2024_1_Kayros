package rest

import (
	"context"

	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
	GetByFilter(ctx context.Context, id alias.CategoryId) ([]*entity.Restaurant, error)
	GetCategoryList(ctx context.Context) ([]*entity.Category, error)
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
	return FromGrpcStructToRestaurantArray(rests), nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	r, err := uc.grpcClient.GetById(ctx, &rest.RestId{Id: uint64(restId)})
	if err != nil {
		return nil, err
	}
	return FromGrpcStructToRestaurant(r), nil
}

func (uc *UsecaseLayer) GetByFilter(ctx context.Context, id alias.CategoryId) ([]*entity.Restaurant, error) {
	rests, err := uc.grpcClient.GetByFilter(ctx, &rest.Id{Id: uint64(id)})
	if err != nil {
		return nil, err
	}
	return FromGrpcStructToRestaurantArray(rests), nil
}

func (uc *UsecaseLayer) GetCategoryList(ctx context.Context) ([]*entity.Category, error) {
	cats, err := uc.grpcClient.GetCategoryList(ctx, nil)
	if err != nil {
		return nil, err
	}
	return FromGrpcStructToCategoryArray(cats), nil
}
