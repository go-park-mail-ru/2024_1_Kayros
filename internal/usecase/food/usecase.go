package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Category, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
}

type UsecaseLayer struct {
	repoFood food.Repo
}

func NewUsecaseLayer(repoFoodProps food.Repo) Usecase {
	return &UsecaseLayer{
		repoFood: repoFoodProps,
	}
}

func (uc *UsecaseLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Category, error) {

	dishes, err := uc.repoFood.GetByRestId(ctx, restId)
	if err != nil {
		return nil, err
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
	return categories, nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	dish, err := uc.repoFood.GetById(ctx, foodId)
	if err != nil {
		return nil, err
	}
	return dish, nil
}
