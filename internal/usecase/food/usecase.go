package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/utils/alias"
)

<<<<<<< HEAD
type Usecase interface {
	GetByRest(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
	AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error)
	UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error)
	DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error)
}

type UsecaseLayer struct {
	repoFood food.Repo
}

func NewUsecaseLayer(repoFoodProps food.Repo) Usecase {
	return &UsecaseLayer{repoFood: repoFoodProps}
}

func (uc *UsecaseLayer) GetByRest(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	return uc.repoFood.GetByRestId(ctx, restId)
}

func (uc *UsecaseLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	return uc.repoFood.GetById(ctx, foodId)
}

func (uc *UsecaseLayer) AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	return uc.repoFood.AddToOrder(ctx, foodId, orderId)
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error) {
	return uc.repoFood.UpdateCountInOrder(ctx, foodId, orderId, count)
}
func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	return uc.repoFood.DeleteFromOrder(ctx, foodId, orderId)
}
=======
type UseCaseInterface interface {
	GetByRest(ctx context.Context, restId uint64) ([]*entity.Food, error)
	GetById(ctx context.Context, id uint64) (*entity.Food, error)
}

type UseCase struct {
	repo food.RepoInterface
}

func NewUseCase(r food.RepoInterface) UseCaseInterface {
	return &UseCase{repo: r}
}

func (uc *UseCase) GetByRest(ctx context.Context, restId uint64) ([]*entity.Food, error) {
	dishes, err := uc.repo.GetByRestId(ctx, restId)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func (uc *UseCase) GetById(ctx context.Context, id uint64) (*entity.Food, error) {
	dish, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dish, nil
}
>>>>>>> 413f5b421db12a295cbeea451991559a66aa908b
