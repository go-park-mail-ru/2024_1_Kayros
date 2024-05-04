package usecase

import (
	"context"
	"sync"

	"2024_1_kayros/microservices/restaurants/internal/repo"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

type Rest interface {
	GetAll(ctx context.Context, _ *rest.Empty) (*rest.RestList, error)
	GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error)
	GetByFilter(ctx context.Context, filter *rest.Filter) (*rest.RestList, error)
	GetCategoryList(ctx context.Context, _ *rest.Empty) (*rest.CategoryList, error)
}

type RestLayer struct {
	rest.UnimplementedRestWorkerServer

	mu       *sync.RWMutex
	repoRest repo.Rest
}

func NewRestLayer(repoRestProps repo.Rest) *RestLayer {
	return &RestLayer{
		UnimplementedRestWorkerServer: rest.UnimplementedRestWorkerServer{},
		repoRest:                      repoRestProps,
		mu:                            &sync.RWMutex{},
	}
}

func (uc *RestLayer) GetAll(ctx context.Context, _ *rest.Empty) (*rest.RestList, error) {
	rests, err := uc.repoRest.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return rests, nil
}

func (uc *RestLayer) GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error) {
	r, err := uc.repoRest.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (uc *RestLayer) GetByFilter(ctx context.Context, filter *rest.Filter) (*rest.RestList, error) {
	rests, err := uc.repoRest.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return rests, nil
}

func (uc *RestLayer) GetCategoryList(ctx context.Context, _ *rest.Empty) (*rest.CategoryList, error) {
	cats, err := uc.repoRest.GetCategoryList(ctx)
	if err != nil {
		return nil, err
	}
	return cats, nil
}
