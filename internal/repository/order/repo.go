package order

import (
	"context"
	"database/sql"
	"errors"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	orderStatus "2024_1_kayros/internal/utils/constants"
)

type Repo interface {
	Create(ctx context.Context, userId alias.UserId, dateOrder string) (*entity.Order, error)
	GetOrdersByUserId(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string, status string) ([]*entity.Order, error)
	GetOrderIdByUserId(ctx context.Context, userId alias.UserId, status string) (alias.OrderId, error)
	GetOrderIdByUserEmail(ctx context.Context, email string, status string) (alias.OrderId, error)
	GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) error
	AddToOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{db: dbProps}
}

// нужно еще добавить проверку на то, есть ли такая запись в таблице уже
func (repo *RepoLayer) Create(ctx context.Context, userId alias.UserId, dateOrder string) (*entity.Order, error) {
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO "Order" (user_id, date_order, status) VALUES ($1, $2, $3)`, uint64(userId), dateOrder, orderStatus.Draft)
	if err != nil {
		return nil, err
	}

	countRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if countRows == 0 {
		return nil, errors.New("Заказ не был добавлен")
	}

	orders, err := repo.GetOrdersByUserId(ctx, userId, orderStatus.Draft)
	if err != nil {
		return nil, err
	}
	return orders[0], nil
}

func (repo *RepoLayer) GetOrdersByUserId(ctx context.Context, userId alias.UserId, status string) ([]*entity.Order, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, user_id, date_order, date_receiving, status, address, 
       				extra_address, sum FROM "Order" WHERE user_id= $1 AND status=$2`, uint64(userId), status)
	if err != nil {
		return nil, err
	}

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
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
		orders = append(orders, &order)
	}

	return orders, nil
}

func (repo *RepoLayer) GetOrdersByUserEmail(ctx context.Context, email string, status string) ([]*entity.Order, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, user_id, date_order, date_receiving, status, address, 
       				extra_address, sum FROM "Order" WHERE email= $1 AND status=$2`, email, status)
	if err != nil {
		return nil, err
	}

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
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
		orders = append(orders, &order)
	}
	return orders, nil
}

func (repo *RepoLayer) GetOrderIdByUserId(ctx context.Context, userId alias.UserId, status string) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id FROM "Order" WHERE user_id= $1 AND status=$2`, uint64(userId), orderStatus.Draft)
	var orderId uint64
	err := row.Scan(&orderId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return alias.OrderId(orderId), nil
}

func (repo *RepoLayer) GetOrderIdByUserEmail(ctx context.Context, email string, status string) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id FROM "Order" WHERE email= $1 AND status=$2`, email, orderStatus.Draft)
	var orderId uint64
	err := row.Scan(&orderId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return alias.OrderId(orderId), nil
}

func (repo *RepoLayer) GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT f.id, f.name, f.weight, f.price, fo.count, f.img_url
				FROM "FoodOrder" AS fo
				JOIN "Food" AS f ON fo.food_id = f.id
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

func (repo *RepoLayer) Update(ctx context.Context, order *entity.Order) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE "Order" SET date_receiving=$1, status=$2, address=$3, extra_address=$4, sum=$5 
               WHERE order_id=$4`, order.DateReceiving, order.Status, order.Address, order.ExtraAddress, order.Sum)
	if err != nil {
		return err
	}

	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return errors.New("Данные о заказе не были обновлены")
	}
	return nil
}

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
		return errors.New("Статус заказа не был обновлён")
	}
	return nil
}

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
		return errors.New("Блюдо не было добавлено в заказ")
	}
	return err
}

func (repo *RepoLayer) UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE "FoodOrder" SET count=$1 WHERE order_id=$2 AND food_id=$3`, count, uint64(orderId), uint64(foodId))
	if err != nil {
		return err
	}

	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return errors.New("Блюдо не было добавлено в заказ")
	}
	return err
}

func (repo *RepoLayer) DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error {
	res, err := repo.db.ExecContext(ctx,
		`DELETE FROM "FoodOrder" WHERE order_id=$1 AND food_id=$2`, uint64(orderId), uint64(foodId))
	if err != nil {
		return err
	}

	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return errors.New("Блюдо не было удалено из заказа")
	}
	return err
}
