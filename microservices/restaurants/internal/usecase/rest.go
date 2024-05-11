package usecase

import (
	"context"
	"errors"
	"sync"

	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/microservices/restaurants/internal/repo"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type Rest interface {
	GetAll(ctx context.Context, _ *rest.Empty) (*rest.RestList, error)
	GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error)
	GetByFilter(ctx context.Context, filter *rest.Id) (*rest.RestList, error)
	GetCategoryList(ctx context.Context, _ *rest.Empty) (*rest.CategoryList, error)
}

type RestLayer struct {
	rest.UnimplementedRestWorkerServer
	mu       *sync.RWMutex
	repoRest repo.Rest
	logger   *zap.Logger
}

func NewRestLayer(repoRestProps repo.Rest, loggerProps *zap.Logger) *RestLayer {
	return &RestLayer{
		UnimplementedRestWorkerServer: rest.UnimplementedRestWorkerServer{},
		repoRest:                      repoRestProps,
		mu:                            &sync.RWMutex{},
		logger:                        loggerProps,
	}
}

func (uc *RestLayer) GetAll(ctx context.Context, _ *rest.Empty) (*rest.RestList, error) {
	rests, err := uc.repoRest.GetAll(ctx)
	if err != nil {
		uc.logger.Error(err.Error())
		return &rest.RestList{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return rests, nil
}

func (uc *RestLayer) GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error) {
	r, err := uc.repoRest.GetById(ctx, id)
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsRestaurantRelation) {
			return &rest.Rest{}, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return &rest.Rest{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return r, nil
}

func (uc *RestLayer) GetByFilter(ctx context.Context, filter *rest.Id) (*rest.RestList, error) {
	rests, err := uc.repoRest.GetByFilter(ctx, filter)
	if err != nil {
		uc.logger.Error(err.Error())
		return &rest.RestList{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return rests, nil
}

func (uc *RestLayer) GetCategoryList(ctx context.Context, _ *rest.Empty) (*rest.CategoryList, error) {
	cats, err := uc.repoRest.GetCategoryList(ctx)
	if err != nil {
		uc.logger.Error(err.Error())
		return &rest.CategoryList{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return cats, nil
}
