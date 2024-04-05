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
	NoBasketError        = "У пользователя нет корзины"
	NotUpdateError       = "Данные о заказе не были обновлены"
	NotUpdateStatusError = "Статус заказа не был обновлен"
	NotAddFood           = "Блюдо не добавлено в заказ"
	NotDeleteFood        = "Блюдо не удалено из заказа"
)

type Repo interface {
	Create(ctx context.Context, userId alias.UserId, dateOrder string) (alias.OrderId, error)
	GetOrders(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error)
	GetBasketId(ctx context.Context, userId alias.UserId) (alias.OrderId, error)
	GetOrderById(ctx context.Context, orderId alias.OrderId) (*entity.Order, error)
	GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error)
	Update(ctx context.Context, order *entity.Order) (alias.OrderId, error)
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error
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
		`INSERT INTO Order (user_id, date_order, status) VALUES ($1, $2, $3)`, uint64(userId), dateOrder, orderStatus.Draft)
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
	rows, err := repo.db.QueryContext(ctx, `SELECT id, user_id, date_order, date_receiving, status, address, 
       				extra_address, sum FROM Order WHERE user_id= $1 AND status=$2`, uint64(userId), status)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Order
	for rows.Next() {
		var order *entity.Order
		err = rows.Scan(&order.Id, &order.UserId, &order.DateOrder, &order.DateReceiving, &order.Status, &order.Address,
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
		orders = append(orders, order)
	}
	return orders, nil
}

// ok
func (repo *RepoLayer) GetOrderById(ctx context.Context, orderId alias.OrderId) (*entity.Order, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id, user_id, date_order, date_receiving, status, address, 
       				extra_address, sum FROM Order WHERE id= $1`, uint64(orderId))
	var order *entity.Order
	err := row.Scan(&order.Id, &order.UserId, &order.DateOrder, &order.DateReceiving, &order.Status, &order.Address,
		&order.ExtraAddress, &order.Sum)
	if err != nil {

		return nil, err
	}
	foodArray, err := repo.GetFood(ctx, orderId)
	order.Food = foodArray
	return order, nil
}

// ok
func (repo *RepoLayer) GetBasketId(ctx context.Context, userId alias.UserId) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT id FROM Order WHERE user_id= $1 AND status=$2", uint64(userId), orderStatus.Draft)
	var orderId uint64
	err := row.Scan(&orderId)
	if errors.Is(err, sql.ErrNoRows) {
		//return 0, fmt.Errorf(NoBasketError)
		return 0, nil
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
				FROM FoodOrder AS fo
				JOIN Food AS f ON fo.food_id = f.id
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
		foodArray = append(foodArray, &food)
	}
	return foodArray, nil
}

// ok
func (repo *RepoLayer) Update(ctx context.Context, order *entity.Order) (alias.OrderId, error) {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE Order SET date_receiving=$1, address=$3, extra_address=$4, sum=$5 
               WHERE order_id=$4`, order.DateReceiving, order.Address, order.ExtraAddress, order.Sum)
	if err != nil {
		return 0, err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if countRows == 0 {
		return 0, fmt.Errorf(NotUpdateError)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return alias.OrderId(id), err
}

// ok
func (repo *RepoLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error {
	res, err := repo.db.ExecContext(ctx, `UPDATE "Order" SET status=$1 WHERE order_id=$2`, status, uint64(orderId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return fmt.Errorf(NotUpdateStatusError)
	}
	return nil
}

// ok
func (repo *RepoLayer) AddToOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO "FoodOrder" (order_id, food_id, count) VALUES ($1, $2, $3)`, uint64(orderId), uint64(foodId), count)
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

	var sum uint64
	row := repo.db.QueryRowContext(ctx,
		`SELECT sum FROM Order WHERE order_id=$1`, uint64(orderId))
	if err != nil {
		return err
	}
	err = row.Scan(&sum)
	if err != nil {
		return err
	}

	var price uint64
	row = repo.db.QueryRowContext(ctx,
		`SELECT price FROM Food WHERE food_id=$1`, uint64(foodId))
	if err != nil {
		return err
	}
	err = row.Scan(&price)
	if err != nil {
		return err
	}

	res, err = repo.db.ExecContext(ctx,
		`UPDATE Order SET sum=$1`, sum+price)
	if err != nil {
		return err
	}
	countRows, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return fmt.Errorf(NotAddFood)
	}
	return nil
}

// ok
func (repo *RepoLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE FoodOrder SET count=$1 WHERE order_id=$2 AND food_id=$3`, count, uint64(orderId), uint64(foodId))
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
	return err
}

// ok
func (repo *RepoLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error {
	res, err := repo.db.ExecContext(ctx,
		`DELETE FROM FoodOrder WHERE order_id=$1 AND food_id=$2`, uint64(orderId), uint64(foodId))
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
	return err
}
