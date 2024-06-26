package usecase

import (
	foodproto "2024_1_kayros/gen/go/food"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/microservices/food/internal/repo"
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type Usecase interface {
	foodproto.UnsafeFoodManagerServer
	GetByRestId(ctx context.Context, restId *foodproto.RestId) (*foodproto.RestCategories, error)
	GetById(ctx context.Context, foodId *foodproto.FoodId) (*foodproto.Food, error)
}

type Layer struct {
	foodproto.UnsafeFoodManagerServer
	repoFood repo.Repo
	logger   *zap.Logger
}

func NewLayer(repoFoodProps repo.Repo, loggerProps *zap.Logger) Usecase {
	return &Layer{
		repoFood: repoFoodProps,
		logger:   loggerProps,
	}
}

func (uc *Layer) GetByRestId(ctx context.Context, restId *foodproto.RestId) (*foodproto.RestCategories, error) {
	dishes, err := uc.repoFood.GetByRestId(ctx, alias.RestId(restId.GetId()))
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	categories := []*entity.Category{}
	if len(dishes) > 0 {
		id := 0
		category := &entity.Category{
			Id:   alias.CategoryId(id),
			Name: dishes[0].Category,
			Food: []*entity.Food{dishes[0]},
		}
		for i := 1; i < len(dishes); i++ {
			if dishes[i].Category != dishes[i-1].Category {
				id++
				categories = append(categories, category)
				category = &entity.Category{
					Id:   alias.CategoryId(id),
					Name: dishes[i].Category,
					Food: []*entity.Food{},
				}
			}
			category.Food = append(category.Food, dishes[i])
		}
	}
	return entity.CnvEntityCtgIntoProtoCtgs(categories), nil
}

func (uc *Layer) GetById(ctx context.Context, foodId *foodproto.FoodId) (*foodproto.Food, error) {
	dish, err := uc.repoFood.GetById(ctx, alias.FoodId(foodId.GetId()))
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsFoodRelation) {
			return nil, grpcerr.NewError(codes.NotFound, myerrors.SqlNoRowsFoodRelation.Error())
		}
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	return entity.CnvEntityFoodIntoProtoFood(dish), nil
}
