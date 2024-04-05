package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/utils/alias"
	"go.uber.org/zap"
)

type Usecase interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
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

func (uc *UsecaseLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	return uc.repoFood.GetByRestId(ctx, restId)
}

func (uc *UsecaseLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	return uc.repoFood.GetById(ctx, foodId)
}
