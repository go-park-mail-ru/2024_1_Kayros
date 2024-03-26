package restaurants

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/restaurants"
)

type RestaurantUseCaseInterface interface {
	GetRestaurants(ctx context.Context) ([]*entity.Restaurant, error)
	GetRestaurantById(ctx context.Context, id int) (*entity.Restaurant, error)
}

type RestaurantUseCase struct {
	repo restaurants.RestaurantRepo
}

func NewRestaurantUseCase(r restaurants.RestaurantRepo) *RestaurantUseCase {
	return &RestaurantUseCase{repo: r}
}

func (uc *RestaurantUseCase) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	var rests []*entity.Restaurant
	rests, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return rests, err
}
func (uc *RestaurantUseCase) GetById(ctx context.Context, id int) (*entity.Restaurant, error) {
	var rest *entity.Restaurant
	rest, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return rest, err
}
