package order

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/order"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

type Usecase interface {
	GetBasketId(ctx context.Context, email string) (alias.OrderId, error)
	GetBasket(ctx context.Context, email string) (*entity.Order, error)
	Create(ctx context.Context, email string) (alias.OrderId, error)
	Update(ctx context.Context, order *entity.Order) (*entity.Order, error)
	UpdateStatus(ctx context.Context, orderId alias.OrderId, currentStatus string, newStatus string) error
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

// ok
func (uc *UsecaseLayer) GetBasketId(ctx context.Context, email string) (alias.OrderId, error) {
	methodName := constants.NameMethodGetBasketId
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	User, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	id, err := uc.repoOrder.GetBasketId(ctx, alias.UserId(User.Id))
	if err != nil {
		return 0, err
	}
	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
	return id, nil
}

// ok
func (uc *UsecaseLayer) GetBasket(ctx context.Context, email string) (*entity.Order, error) {
	methodName := constants.NameMethodGetBasket
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	User, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil {
		return nil, err
	}
	orders, err := uc.repoOrder.GetOrders(ctx, alias.UserId(User.Id), constants.Draft)
	if err != nil {
		return nil, err
	}
	return orders[0], nil
}

// ok
func (uc *UsecaseLayer) Create(ctx context.Context, email string) (alias.OrderId, error) {
	methodName := constants.NameMethodCreateOrder
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	User, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil {
		return 0, err
	}
	dateOrder := time.Now()
	dateOrderForDB := dateOrder.Format("2024-04-02 23:34:54")
	id, err := uc.repoOrder.Create(ctx, alias.UserId(User.Id), dateOrderForDB)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (uc *UsecaseLayer) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	id, err := uc.repoOrder.Update(ctx, order)
	if err != nil {
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	return updatedOrder, nil
}

func (uc *UsecaseLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, currentStatus string, newStatus string) error {
	if currentStatus == constants.Draft && newStatus == constants.Payed {
		return uc.repoOrder.UpdateStatus(ctx, orderId, newStatus)
	}
	return fmt.Errorf("Статус заказа не был обновлен")
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
