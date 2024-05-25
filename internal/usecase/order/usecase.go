package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/repository/food"
	"2024_1_kayros/internal/repository/order"

	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/gen/go/user"

	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"

	"google.golang.org/grpc/codes"
)

type Usecase interface {
	GetBasketId(ctx context.Context, email string) (alias.OrderId, error)
	GetBasketIdNoAuth(ctx context.Context, token string) (alias.OrderId, error)
	GetBasket(ctx context.Context, email string) (*entity.Order, error)
	GetBasketNoAuth(ctx context.Context, token string) (*entity.Order, error)
	GetOrderById(ctx context.Context, id alias.OrderId) (*entity.Order, error)
	Create(ctx context.Context, email string) (alias.OrderId, error)
	GetCurrentOrders(ctx context.Context, email string) ([]*entity.ShortOrder, error)
	GetArchiveOrders(ctx context.Context, email string) ([]*entity.ShortOrder, error)
	CreateNoAuth(ctx context.Context, token string) (alias.OrderId, error)
	UpdateAddress(ctx context.Context, FullAddress dto.FullAddress, orderId alias.OrderId) error
	Pay(ctx context.Context, orderId alias.OrderId, currentStatus string, email string, userId alias.UserId) (*entity.Order, error)
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) (*entity.Order, error)
	Clean(ctx context.Context, orderId alias.OrderId) error
	AddFoodToOrder(ctx context.Context, foodId alias.FoodId, count uint32, orderId alias.OrderId) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) (*entity.Order, error)
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) (*entity.Order, error)
	UpdateSum(ctx context.Context, sum uint64, orderId alias.OrderId) error

	CheckPromocode(ctx context.Context, email string, codeName string, basketId alias.OrderId) (*entity.Promocode, error)
	SetPromocode(ctx context.Context, orderId alias.OrderId, code *entity.Promocode) (uint64, error)
	GetPromocodeByOrder(ctx context.Context, orderId *alias.OrderId) (*entity.Promocode, error)
	DeletePromocode(ctx context.Context, orderId alias.OrderId) error
}

type UsecaseLayer struct {
	repoOrder      order.Repo
	userGrpcClient user.UserManagerClient
	repoFood       food.Repo
	restGrpcClient rest.RestWorkerClient
}

func NewUsecaseLayer(repoOrderProps order.Repo, repoFoodProps food.Repo, repoUserProps user.UserManagerClient, repoRestProps rest.RestWorkerClient) Usecase {
	return &UsecaseLayer{
		repoOrder:      repoOrderProps,
		userGrpcClient: repoUserProps,
		repoFood:       repoFoodProps,
		restGrpcClient: repoRestProps,
	}
}

func (uc *UsecaseLayer) GetBasketId(ctx context.Context, email string) (alias.OrderId, error) {
	u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return 0, myerrors.SqlNoRowsUserRelation
		}
		return 0, err
	}
	id, err := uc.repoOrder.GetBasketId(ctx, alias.UserId(u.Id))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (uc *UsecaseLayer) GetBasketIdNoAuth(ctx context.Context, token string) (alias.OrderId, error) {
	id, err := uc.repoOrder.GetBasketIdNoAuth(ctx, token)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (uc *UsecaseLayer) GetBasket(ctx context.Context, email string) (*entity.Order, error) {
	u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return nil, err
	}
	orders, err := uc.repoOrder.GetOrders(ctx, alias.UserId(u.Id), constants.Draft)
	if err != nil {
		return nil, err
	}
	basket := orders[0]
	if len(basket.Food) != 0 {
		basket.RestaurantId = basket.Food[0].RestaurantId
	}
	return basket, nil
}

func (uc *UsecaseLayer) GetBasketNoAuth(ctx context.Context, token string) (*entity.Order, error) {
	Order, err := uc.repoOrder.GetBasketNoAuth(ctx, token)
	if err != nil {
		return nil, err
	}
	if len(Order.Food) != 0 {
		Order.RestaurantId = Order.Food[0].RestaurantId
	}
	return Order, nil
}

func (uc *UsecaseLayer) GetOrderById(ctx context.Context, id alias.OrderId) (*entity.Order, error) {
	Order, err := uc.repoOrder.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(Order.Food) != 0 {
		Order.RestaurantId = Order.Food[0].RestaurantId
		r, err := uc.restGrpcClient.GetById(ctx, &rest.RestId{Id: Order.RestaurantId})
		if err != nil {
			if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsRestaurantRelation) {
				return nil, myerrors.SqlNoRowsRestaurantRelation
			}
			return nil, err
		}
		Order.RestaurantName = r.Name
	}
	return Order, nil
}

func (uc *UsecaseLayer) GetCurrentOrders(ctx context.Context, email string) ([]*entity.ShortOrder, error) {
	u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return nil, err
	}
	orders, err := uc.repoOrder.GetOrders(ctx, alias.UserId(u.Id), constants.Payed, constants.Created, constants.Cooking, constants.OnTheWay)
	if err != nil {
		return nil, err
	}
	res := []*entity.ShortOrder{}
	for _, o := range orders {
		id := o.Food[0].RestaurantId
		rest, err := uc.restGrpcClient.GetById(ctx, &rest.RestId{Id: id})
		if err != nil {
			if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsRestaurantRelation) {
				return nil, myerrors.SqlNoRowsRestaurantRelation
			}
			return nil, err
		}
		res = append(res, &entity.ShortOrder{
			Id:             o.Id,
			UserId:         o.UserId,
			Status:         o.Status,
			Time:           "",
			RestaurantId:   o.RestaurantId,
			RestaurantName: rest.Name,
		})
	}
	return res, nil
}

func (uc *UsecaseLayer) GetArchiveOrders(ctx context.Context, email string) ([]*entity.ShortOrder, error) {
	u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return nil, err
	}
	orders, err := uc.repoOrder.GetOrders(ctx, alias.UserId(u.Id), constants.Delivered)
	if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
		return []*entity.ShortOrder{}, nil
	}
	if err != nil {
		return nil, err
	}
	res := []*entity.ShortOrder{}
	for _, o := range orders {
		id := o.Food[0].RestaurantId
		rest, err := uc.restGrpcClient.GetById(ctx, &rest.RestId{Id: id})
		if err != nil {
			if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsRestaurantRelation) {
				return nil, myerrors.SqlNoRowsRestaurantRelation
			}
			return nil, err
		}
		res = append(res, &entity.ShortOrder{
			Id:             o.Id,
			UserId:         o.UserId,
			Status:         o.Status,
			Time:           o.OrderCreatedAt,
			RestaurantId:   o.RestaurantId,
			RestaurantName: rest.Name,
			Sum:            uint32(o.Sum),
		})
	}
	return res, nil
}

func (uc *UsecaseLayer) Create(ctx context.Context, email string) (alias.OrderId, error) {
	u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return 0, myerrors.SqlNoRowsUserRelation
		}
		return 0, err
	}
	id, err := uc.repoOrder.Create(ctx, alias.UserId(u.Id))
	if err != nil {
		return 0, err
	}
	return id, err
}

func (uc *UsecaseLayer) CreateNoAuth(ctx context.Context, token string) (alias.OrderId, error) {
	id, err := uc.repoOrder.CreateNoAuth(ctx, token)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, FullAddress dto.FullAddress, orderId alias.OrderId) error {
	_, err := uc.repoOrder.UpdateAddress(ctx, FullAddress.Address, FullAddress.ExtraAddress, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsecaseLayer) Pay(ctx context.Context, orderId alias.OrderId, currentStatus string, email string, userId alias.UserId) (*entity.Order, error) {
	if currentStatus != constants.Draft {
		return nil, myerrors.AlreadyPayed
	}
	if userId == 0 {
		u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
		if err != nil {
			if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
				return nil, myerrors.SqlNoRowsUserRelation
			}
			return nil, err
		}
		err = uc.repoOrder.SetUser(ctx, orderId, alias.UserId(u.Id))
		if err != nil {
			return nil, err
		}
	}
	id, err := uc.repoOrder.UpdateStatus(ctx, orderId, constants.Payed)
	if err != nil {
		return nil, err
	}
	Order, err := uc.repoOrder.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(Order.Food) != 0 {
		Order.RestaurantId = Order.Food[0].RestaurantId
	}
	return Order, nil
}

func (uc *UsecaseLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) (*entity.Order, error) {
	id, err := uc.repoOrder.UpdateStatus(ctx, orderId, status)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	Order, err := uc.repoOrder.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(Order.Food) != 0 {
		Order.RestaurantId = Order.Food[0].RestaurantId
	}
	return Order, nil
}

func (uc *UsecaseLayer) Clean(ctx context.Context, orderId alias.OrderId) error {
	err := uc.repoOrder.DeleteBasket(ctx, orderId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (uc *UsecaseLayer) AddFoodToOrder(ctx context.Context, foodId alias.FoodId, count uint32, orderId alias.OrderId) error {
	//err := uc.repoOrder.AddToOrder(ctx, orderId, foodId, count)
	//if err != nil {
	//	return err
	//}
	//return err
	//получаем блюдо по id
	inputFood, err := uc.repoFood.GetById(ctx, foodId)
	if err != nil {
		return err
	}
	//получаем заказ по id
	Order, err := uc.repoOrder.GetOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	//fmt.Println(Order.Food[0].RestaurantId)
	//если ресторан блюд в корзине не совпадает с рестораном откуда новое блюдо
	//то чистим корзину
	if len(Order.Food) > 0 && (inputFood.RestaurantId != Order.Food[0].RestaurantId) {
		err = uc.repoOrder.CleanBasket(ctx, orderId)
		if err != nil {
			return err
		}
	}
	err = uc.repoOrder.AddToOrder(ctx, orderId, foodId, count)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsecaseLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) (*entity.Order, error) {
	err := uc.repoOrder.UpdateCountInOrder(ctx, orderId, foodId, count)
	if err != nil {
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}
	if len(updatedOrder.Food) != 0 {
		updatedOrder.RestaurantId = updatedOrder.Food[0].RestaurantId
	}
	return updatedOrder, nil
}

func (uc *UsecaseLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) (*entity.Order, error) {
	err := uc.repoOrder.DeleteFromOrder(ctx, orderId, foodId)
	if err != nil {
		return nil, err
	}
	updatedOrder, err := uc.repoOrder.GetOrderById(ctx, orderId)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			return nil, myerrors.SuccessCleanRu
		}
		return nil, err
	}
	if len(updatedOrder.Food) != 0 {
		updatedOrder.RestaurantId = updatedOrder.Food[0].RestaurantId
	}
	return updatedOrder, nil
}

// проверяет 4 вида промокодов
func (uc *UsecaseLayer) CheckPromocode(ctx context.Context, email string, codeName string, basketId alias.OrderId) (*entity.Promocode, error) {
	code, err := uc.repoOrder.GetPromocode(ctx, codeName)
	fmt.Println(code, err)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsPromocodeRelation) {
			return nil, nil
		}
		return nil, err
	}
	if code == nil {
		return nil, nil
	}
	// layout := "2024-05-28 16:52:48+00:00"
	// date, err := time.Parse(layout, code.Date)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	date := code.Date
	if date.Before(time.Now()) {
		return nil, myerrors.OverDatePromocode
	}
	//берем юзера, так как дальше нужно проверять его заказы
	u, err := uc.userGrpcClient.GetData(ctx, &user.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return nil, err
	}

	//промокод на первый заказ в сервисе
	if code.Type == "first" {
		//тут немного другую функцию вызывать
		count, err := uc.repoOrder.OrdersCount(ctx, alias.UserId(u.Id), constants.Delivered)
		//err :=
		if err != nil && !errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			return nil, err
		}
		if count > 1 {
			return nil, myerrors.OncePromocode
		}
	} else if code.Type == "sum" { //промокод от определенной суммы
		sum, err := uc.repoOrder.GetOrderSum(ctx, basketId)
		if err != nil {
			return nil, err
		}
		if uint64(sum) < code.Sum {
			return nil, myerrors.SumPromocode
		}
	}
	//промокод, который можно применить один раз
	if code.Type == "once" {
		//тут немного другую функцию вызывать
		err := uc.repoOrder.WasPromocodeUsed(ctx, alias.UserId(u.Id), code.Id)
		fmt.Println("wasCodeUsed", err)
		if err != nil {
			return nil, err
		}
	} else if code.Type == "rest" { //промокод на первый заказ в рестике
		//тут немного другую функцию вызывать
		orders, err := uc.GetArchiveOrders(ctx, email)
		if err != nil {
			return nil, err
		}
		for _, o := range orders {
			if o.RestaurantId == code.Rest {
				//сюда попадут только те заказы, которые были сделаны в ресторане,
				//который относится к введеному промокоду
				err = uc.repoOrder.WasRestPromocodeUsed(ctx, alias.OrderId(o.Id), code.Id)
				//если промокод был введен, то выкидываемся из функции
				//err=nil, когда не нашлось записей, то есть промокод не был применен
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return code, nil
}

func (uc *UsecaseLayer) SetPromocode(ctx context.Context, orderId alias.OrderId, code *entity.Promocode) (uint64, error) {
	sum, err := uc.repoOrder.SetPromocode(ctx, orderId, code.Id)
	if err != nil {
		return 0, err
	}
	sum = sum * (100 - uint64(code.Sale)) / 100
	return sum, nil
}

func (uc *UsecaseLayer) GetPromocodeByOrder(ctx context.Context, orderId *alias.OrderId) (*entity.Promocode, error) {
	code, err := uc.repoOrder.GetPromocodeByOrder(ctx, orderId)
	fmt.Println(code, err)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsPromocodeRelation) {
			return nil, nil
		}
		return nil, err
	}
	return code, nil
}

func (uc *UsecaseLayer) DeletePromocode(ctx context.Context, orderId alias.OrderId) error {
	err := uc.repoOrder.DeletePromocode(ctx, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsecaseLayer) UpdateSum(ctx context.Context, sum uint64, orderId alias.OrderId) error {
	err := uc.repoOrder.UpdateSum(ctx, uint32(sum), orderId)
	if err != nil {
		return err
	}
	return nil
}
