package food

import (
	"context"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

type Usecase interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Category, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
}

type UsecaseLayer struct {
	repoFood food.Repo
	logger   *zap.Logger
}

func NewUsecaseLayer(repoFoodProps food.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoFood: repoFoodProps,
		logger:   loggerProps,
	}
}

func (uc *UsecaseLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Category, error) {
	methodName := constants.NameMethodGetFoodByRest
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	dishes, err := uc.repoFood.GetByRestId(ctx, requestId, restId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	categories := []*entity.Category{}
	if len(dishes) > 0 {
		id := 1
		category := &entity.Category{
			Id:   alias.CategoryId(id),
			Name: dishes[0].Category,
			Food: []*entity.Food{dishes[0]},
		}
		categories = append(categories, category)
		for i := 1; i < len(dishes); i++ {
			if dishes[i].Category != dishes[i-1].Category {
				id++
				category = &entity.Category{
					Id:   alias.CategoryId(id),
					Name: dishes[i].Category,
					Food: []*entity.Food{},
				}
				categories = append(categories, category)
			}
			category.Food = append(category.Food, dishes[i])
		}
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return categories, nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	methodName := constants.NameMethodGetFoodById
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	dish, err := uc.repoFood.GetById(ctx, requestId, foodId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return dish, nil
}
