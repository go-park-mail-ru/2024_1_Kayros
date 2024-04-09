package food

import (
	"context"
	"fmt"

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
	var categories []*entity.Category
	category := &entity.Category{
		Name: dishes[0].Category,
		Food: []*entity.Food{dishes[0]},
	}
	for i := 1; i < len(dishes); i++ {
		fmt.Println(i, dishes[i])
		if dishes[i].Category != dishes[i-1].Category {
			fmt.Println(i, "категория сменилась", category)
			categories = append(categories, category)
			for _, v := range categories {
				fmt.Println(v.Name)
			}
			category = &entity.Category{
				Name: dishes[i].Category,
				Food: []*entity.Food{},
			}
		}
		category.Food = append(category.Food, dishes[i])
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
