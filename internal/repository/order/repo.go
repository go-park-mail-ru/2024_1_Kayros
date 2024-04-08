package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	orderStatus "2024_1_kayros/internal/utils/constants"
)

const (
	CreateError          = "Не удалось создать заказ"
	NoBasketError        = "У Вас нет корзины"
	NotUpdateError       = "Данные о заказе не были обновлены"
	NotUpdateStatusError = "Заказ не оплачен"
	NotAddFood           = "Блюдо не добавлено в заказ"
	NotDeleteFood        = "Блюдо не удалено из заказа"
)

type Repo interface {
	Create(ctx context.Context, userId alias.UserId, dateOrder string) (alias.OrderId, error)
	GetOrders(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error)
	GetBasketId(ctx context.Context, userId alias.UserId) (alias.OrderId, error)
	GetOrderById(ctx context.Context, orderId alias.OrderId) (*entity.Order, error)
	GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error)
	UpdateAddress(ctx context.Context, address string, extraAddress string, orderId alias.OrderId) (alias.OrderId, error)
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) (alias.OrderId, error)
	AddToOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error
}

type RepoLayer struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRepoLayer(dbProps *sql.DB, loggerProps *zap.Logger) Repo {
	return &RepoLayer{
		db:     dbProps,
		logger: loggerProps,
	}
}

// ok
func (repo *RepoLayer) Create(ctx context.Context, userId alias.UserId, dateOrder string) (alias.OrderId, error) {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO order (user_id, date_order, status) VALUES ($1, $2, $3)`, uint64(userId), dateOrder, orderStatus.Draft)
	if err != nil {
		return 0, err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if countRows == 0 {
		return 0, fmt.Errorf(CreateError)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return alias.OrderId(id), err
}

// ok
func (repo *RepoLayer) GetOrders(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT id, user_id, created_at, received_at, status, address, 
       				extra_address, sum FROM order WHERE user_id= $1 AND status=$2`, uint64(userId), status)
	fmt.Println("in repo", err)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var orders []*entity.Order
	fmt.Println(userId, status)
	for rows.Next() {
		var order entity.OrderDB
		err = rows.Scan(&order.Id, &order.UserId, &order.CreatedAt, &order.ReceivedAt, &order.Status, &order.Address,
			&order.ExtraAddress, &order.Sum)
		if err != nil {
			return nil, err
		}
		var foodArray []*entity.FoodInOrder
		foodArray, err = repo.GetFood(ctx, alias.OrderId(order.Id))
		if err != nil {
			return nil, err
		}
		order.Food = foodArray
		orders = append(orders, entity.ToOrder(&order))
	}
	return orders, nil
}

// ok
func (repo *RepoLayer) GetOrderById(ctx context.Context, orderId alias.OrderId) (*entity.Order, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id, user_id, created_at, received_at, status, address, 
       				extra_address, sum FROM order WHERE id= $1`, uint64(orderId))
	var order entity.OrderDB
	err := row.Scan(&order.Id, &order.UserId, &order.CreatedAt, &order.ReceivedAt, &order.Status, &order.Address,
		&order.ExtraAddress, &order.Sum)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	foodArray, err := repo.GetFood(ctx, orderId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	order.Food = foodArray
	return entity.ToOrder(&order), nil
}

// ok
func (repo *RepoLayer) GetBasketId(ctx context.Context, userId alias.UserId) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id FROM order WHERE user_id= $1 AND status=$2`, uint64(userId), orderStatus.Draft)
	var orderId uint64
	err := row.Scan(&orderId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf(NoBasketError)
	}
	if err != nil {
		return 0, err
	}
	return alias.OrderId(orderId), nil
}

// ok
func (repo *RepoLayer) GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT f.id, f.name, f.weight, f.price, fo.count, f.img_url
				FROM food_order AS fo
				JOIN food AS f ON fo.food_id = f.id
				WHERE fo.order_id = $1`, uint64(orderId))
	if err != nil {
		return nil, err
	}

	var foodArray []*entity.FoodInOrder
	for rows.Next() {
		var food entity.FoodInOrder
		err = rows.Scan(&food.Id, &food.Name, &food.Weight, &food.Price, &food.Count, &food.ImgUrl)
		if err != nil {
			return nil, err
		}
		fmt.Println(food)
		foodArray = append(foodArray, &food)
	}
	return foodArray, nil
}

// ok
func (repo *RepoLayer) UpdateAddress(ctx context.Context, address string, extraAddress string, orderId alias.OrderId) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx,
		`UPDATE order SET address=$1, extra_address=$2
               WHERE id=$3 RETURNING id`, address, extraAddress, uint64(orderId))
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return alias.OrderId(id), err
}

// ok
func (repo *RepoLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, `UPDATE order SET status=$1 WHERE id=$2 RETURNING id`, status, uint64(orderId))
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return alias.OrderId(id), err
}

func (repo *RepoLayer) GetOrderSum(ctx context.Context, orderId alias.OrderId) (uint32, error) {
	var sum uint32
	row := repo.db.QueryRowContext(ctx,
		`SELECT sum FROM order WHERE id=$1`, uint64(orderId))
	err := row.Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}
func (repo *RepoLayer) GetFoodPrice(ctx context.Context, foodId alias.FoodId) (uint32, error) {
	var price uint32
	row := repo.db.QueryRowContext(ctx,
		`SELECT price FROM food WHERE id=$1`, uint64(foodId))
	err := row.Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (repo *RepoLayer) GetFoodCount(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (uint32, error) {
	var count uint32
	row := repo.db.QueryRowContext(ctx,
		`SELECT count FROM food_order WHERE order_id=$1 AND food_id=$2`, uint64(orderId), uint64(foodId))
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	fmt.Println(count)
	return count, nil
}

func (repo *RepoLayer) UpdateSum(ctx context.Context, sum uint32, orderId alias.OrderId) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE order SET sum=$1 WHERE id=$2`, sum, uint64(orderId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return err
	}
	return nil
}

// ok
func (repo *RepoLayer) AddToOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO food_order (order_id, food_id, count) VALUES ($1, $2, $3)`, uint64(orderId), uint64(foodId), count)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return fmt.Errorf(NotAddFood)
	}

	sum, err := repo.GetOrderSum(ctx, orderId)
	if err != nil {
		return err
	}

	price, err := repo.GetFoodPrice(ctx, foodId)
	if err != nil {
		return err
	}
	sum = sum + count*price
	err = repo.UpdateSum(ctx, sum, orderId)
	if err != nil {
		return err
	}
	return nil
}

// ok
func (repo *RepoLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	currentCount, err := repo.GetFoodCount(ctx, foodId, orderId)
	if err != nil {
		return err
	}
	res, err := repo.db.ExecContext(ctx,
		`UPDATE food_order SET count=$1 WHERE order_id=$2 AND food_id=$3`, count, uint64(orderId), uint64(foodId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return fmt.Errorf(NotAddFood)
	}

	price, err := repo.GetFoodPrice(ctx, foodId)
	if err != nil {
		return err
	}
	sum, err := repo.GetOrderSum(ctx, orderId)
	if err != nil {
		return err
	}
	//cCount - 2, count - 5, сумму надо увеличивать
	if num := int(count) - int(currentCount); num > 0 {
		sum = sum + (count-currentCount)*price
	} else {
		sum = sum - (currentCount-count)*price
	}
	fmt.Println(sum)
	return repo.UpdateSum(ctx, sum, orderId)
}

// ok
func (repo *RepoLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error {
	count, err := repo.GetFoodCount(ctx, foodId, orderId)
	if err != nil {
		return err
	}
	res, err := repo.db.ExecContext(ctx,
		`DELETE FROM food_order WHERE order_id=$1 AND food_id=$2`, uint64(orderId), uint64(foodId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return fmt.Errorf(NotDeleteFood)
	}
	price, err := repo.GetFoodPrice(ctx, foodId)
	if err != nil {
		return err
	}
	sum, err := repo.GetOrderSum(ctx, orderId)
	if err != nil {
		return err
	}
	sum = sum - count*price
	fmt.Println(sum)
	return repo.UpdateSum(ctx, sum, orderId)
}
