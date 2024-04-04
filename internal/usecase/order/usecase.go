package order

import (
	"context"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	"go.uber.org/zap"
)

type Usecase interface {
	GetOrdersByUserEmail(ctx context.Context, email string, status string) ([]*entity.Order, error)
	GetOrdersByUserId(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error)
	GetOrderIdByUserId(ctx context.Context, userId alias.UserId, status string) (alias.OrderId, error)
	GetOrderIdByUserEmail(ctx context.Context, email string, status string) (alias.OrderId, error)
	Create(ctx context.Context, email string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error
	AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error
}

type UsecaseLayer struct {
	repoOrder order.Repo
	repoUser  user.Repo
	logger    *zap.Logger
}

func NewUsecaseLayer(repoOrderProps order.Repo, repoUserProps user.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoOrder: repoOrderProps,
		repoUser:  repoUserProps,
		logger:    loggerProps,
	}
}

func (uc *UsecaseLayer) GetOrdersByUserEmail(ctx context.Context, email string, status string) ([]*entity.Order, error) {
	return uc.repoOrder.GetOrdersByUserEmail(ctx, email, status)
}

func (uc *UsecaseLayer) GetOrdersByUserId(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error) {
	return uc.repoOrder.GetOrdersByUserId(ctx, userId, status)
}

func (uc *UsecaseLayer) GetOrderIdByUserId(ctx context.Context, userId alias.UserId, status string) (alias.OrderId, error) {
	return uc.repoOrder.GetOrderIdByUserId(ctx, userId, status)
}

func (uc *UsecaseLayer) GetOrderIdByUserEmail(ctx context.Context, email string, status string) (alias.OrderId, error) {
	return uc.repoOrder.GetOrderIdByUserEmail(ctx, email, status)
}

func (uc *UsecaseLayer) Create(ctx context.Context, email string) (*entity.Order, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	dateOrder := time.Now()
	dateOrderForDB := dateOrder.Format("2024-04-02 23:34:54")
	ord, err := uc.repoOrder.Create(ctx, alias.UserId(u.Id), dateOrderForDB)
	if err != nil {
		return nil, err
	}
	return ord, err
}

func (uc *UsecaseLayer) Update(ctx context.Context, order *entity.Order) error {

	return uc.repoOrder.Update(ctx, order)
}

func (uc *UsecaseLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error {
	return uc.repoOrder.UpdateStatus(ctx, orderId, status)
}

func (uc *UsecaseLayer) AddFoodToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) error {
	return uc.repoOrder.AddToOrder(ctx, orderId, foodId, 1)
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	return uc.repoOrder.UpdateCountInOrder(ctx, orderId, foodId, count)
}

func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error {
	return uc.repoOrder.DeleteFromOrder(ctx, orderId, foodId)
}
