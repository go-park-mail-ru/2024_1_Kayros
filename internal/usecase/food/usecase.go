package food

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/food"
)

type Usecase interface {
	GetByRest(ctx context.Context, restId int) ([]*entity.Food, error)
	GetById(ctx context.Context, id int) (*entity.Food, error)
}

<<<<<<< HEAD
type UseCase struct {
	repo *food.Repo
}

func NewUseCase(r *food.Repo) *UseCase {
	return &UseCase{repo: r}
=======
type UsecaseLayer struct {
	repo food.Repo
}

func NewUsecase(r food.Repo) Usecase {
	return &UsecaseLayer{repo: r}
>>>>>>> 253e2ec7cad8ed399576ea17ce8101c3e2c8cca5
}

func (uc *UsecaseLayer) GetByRest(ctx context.Context, restId int) ([]*entity.Food, error) {
	var dishes []*entity.Food
	dishes, err := uc.repo.GetByRest(ctx, restId)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func (uc *UsecaseLayer) GetById(ctx context.Context, id int) (*entity.Food, error) {
	var dish *entity.Food
	dish, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dish, nil
}
<<<<<<< HEAD
=======

func (uc *UsecaseLayer) AddToOrder(ctx context.Context, foodId int, orderId int) error {
	err := uc.repo.AddToOrder(ctx, foodId, orderId)
	return err
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, foodId int, orderId int, count int) error {
	err := uc.repo.UpdateCountInOrder(ctx, foodId, orderId, count)
	return err
}
func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, foodId int, orderId int) error {
	err := uc.repo.DeleteFromOrder(ctx, foodId, orderId)
	return err
}
>>>>>>> 253e2ec7cad8ed399576ea17ce8101c3e2c8cca5
