package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
)

type Repo interface {
	Create(ctx context.Context, userId alias.UserId) (alias.OrderId, error)
	CreateNoAuth(ctx context.Context, unauthId string) (alias.OrderId, error)
	GetOrders(ctx context.Context, userId alias.UserId, status ...string) ([]*entity.Order, error)
	GetBasketNoAuth(ctx context.Context, unauthId string) (*entity.Order, error)
	GetBasketId(ctx context.Context, userId alias.UserId) (alias.OrderId, error)
	GetBasketIdNoAuth(ctx context.Context, unauthId string) (alias.OrderId, error)
	GetOrderById(ctx context.Context, orderId alias.OrderId) (*entity.Order, error)
	GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error)
	UpdateAddress(ctx context.Context, address string, extraAddress string, orderId alias.OrderId) (alias.OrderId, error)
	UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) (alias.OrderId, error)
	AddToOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	UpdateCountInOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error
	DeleteFromOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId) error
	CleanBasket(ctx context.Context, orderId alias.OrderId) error
	DeleteBasket(ctx context.Context, orderId alias.OrderId) error
	SetUser(ctx context.Context, orderId alias.OrderId, userId alias.UserId) error
}

type RepoLayer struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{
		db: dbProps,
	}
}

func (repo *RepoLayer) Create(ctx context.Context, userId alias.UserId) (alias.OrderId, error) {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row := repo.db.QueryRowContext(ctx,
		`INSERT INTO "order" (user_id, created_at, updated_at, status) VALUES ($1, $2, $3, $4) RETURNING id`, uint64(userId), timeNow, timeNow, cnst.Draft)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.BasketCreate
		}
		return 0, err
	}
	return alias.OrderId(id), nil
}

func (repo *RepoLayer) CreateNoAuth(ctx context.Context, unauthId string) (alias.OrderId, error) {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row := repo.db.QueryRowContext(ctx,
		`INSERT INTO "order" (unauth_token, created_at, updated_at, status) VALUES ($1, $2, $3, $4) RETURNING id`, unauthId, timeNow, timeNow, cnst.Draft)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.BasketCreate
		}
		return 0, err
	}
	fmt.Println(id)
	return alias.OrderId(id), nil
}

func (repo *RepoLayer) GetOrders(ctx context.Context, userId alias.UserId, status ...string) ([]*entity.Order, error) {
	var rows *sql.Rows
	var err error
	if len(status) == 1 {
		rows, err = repo.db.QueryContext(ctx, `SELECT id, user_id, created_at, status, address,
      				extra_address, sum FROM "order" WHERE user_id= $1 AND status=$2`, uint64(userId), status[0])
	} else {
		str := "$2"
		for i := range len(status) - 1 {
			str = str + ", $" + strconv.Itoa(i+3)
		}
		query := `SELECT id, user_id, created_at, status, address, 
       			extra_address, sum FROM "order" WHERE user_id= $1 AND status IN (` + str + `)`
		args := make([]interface{}, len(status)+1)
		args[0] = uint64(userId)
		for i, a := range status {
			args[i+1] = a
		}
		rows, err = repo.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	orders := []*entity.Order{}
	for rows.Next() {
		var order entity.OrderDB
		err = rows.Scan(&order.Id, &order.UserId, &order.CreatedAt, &order.Status, &order.Address,
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
	if len(orders) == 0 {
		return nil, myerrors.SqlNoRowsOrderRelation
	}
	return orders, nil
}

func (repo *RepoLayer) GetBasketNoAuth(ctx context.Context, unauthId string) (*entity.Order, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id, created_at, updated_at, received_at, status, address, 
       				extra_address, sum FROM "order" WHERE unauth_token= $1 AND status=$2`, unauthId, cnst.Draft)
	var order entity.OrderDB
	err := row.Scan(&order.Id, &order.CreatedAt, &order.UpdatedAt, &order.ReceivedAt, &order.Status, &order.Address,
		&order.ExtraAddress, &order.Sum)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsOrderRelation
		}
		return nil, err
	}
	foodArray, err := repo.GetFood(ctx, alias.OrderId(order.Id))
	if err != nil {
		return nil, err
	}
	order.Food = foodArray
	return entity.ToOrder(&order), nil
}

func (repo *RepoLayer) GetOrderById(ctx context.Context, orderId alias.OrderId) (*entity.Order, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id, user_id, order_created_at, delivered_at, status, address,
      				extra_address, sum FROM "order" WHERE id= $1`, uint64(orderId))
	var order entity.OrderDB
	err := row.Scan(&order.Id, &order.UserId, &order.OrderCreatedAt, &order.DeliveredAt,
		&order.Status, &order.Address, &order.ExtraAddress, &order.Sum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsOrderRelation
		}
		return nil, err
	}
	fmt.Println(order.UserId)
	foodArray, err := repo.GetFood(ctx, orderId)
	if err != nil {
		return nil, err
	}
	order.Food = foodArray
	return entity.ToOrder(&order), nil
}

func (repo *RepoLayer) GetBasketId(ctx context.Context, userId alias.UserId) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id FROM "order" WHERE user_id= $1 AND status=$2`, uint64(userId), cnst.Draft)
	var orderId uint64
	err := row.Scan(&orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsOrderRelation
		}
		return 0, err
	}
	return alias.OrderId(orderId), nil
}

func (repo *RepoLayer) GetBasketIdNoAuth(ctx context.Context, unauthId string) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx, `SELECT id FROM "order" WHERE unauth_token=$1 AND status=$2`, unauthId, cnst.Draft)
	var orderId uint64
	err := row.Scan(&orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsOrderRelation
		}
		return 0, err
	}
	return alias.OrderId(orderId), nil
}

func (repo *RepoLayer) GetFood(ctx context.Context, orderId alias.OrderId) ([]*entity.FoodInOrder, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT f.id, f.name, f.weight, f.price, fo.count, f.img_url, f.restaurant_id
				FROM food_order AS fo
				JOIN food AS f ON fo.food_id = f.id
				WHERE fo.order_id = $1`, uint64(orderId))
	if err != nil {
		return nil, err
	}

	var foodArray []*entity.FoodInOrder
	for rows.Next() {
		var food entity.FoodInOrder
		err = rows.Scan(&food.Id, &food.Name, &food.Weight, &food.Price, &food.Count, &food.ImgUrl, &food.RestaurantId)
		if err != nil {
			return nil, err
		}
		foodArray = append(foodArray, &food)
	}
	return foodArray, nil
}

func (repo *RepoLayer) UpdateAddress(ctx context.Context, address string, extraAddress string, orderId alias.OrderId) (alias.OrderId, error) {
	row := repo.db.QueryRowContext(ctx,
		`UPDATE "order" SET address=$1, extra_address=$2
              WHERE id=$3 RETURNING id`, address, extraAddress, uint64(orderId))
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsOrderRelation
		}
		return 0, err
	}
	return alias.OrderId(id), nil
}

func (repo *RepoLayer) UpdateStatus(ctx context.Context, orderId alias.OrderId, status string) (alias.OrderId, error) {
	var id uint64
	var err error
	if status == cnst.Created {
		timeNow := time.Now().UTC().Format(cnst.Timestamptz)
		err = repo.db.QueryRowContext(ctx, `UPDATE "order" SET status=$1, order_created_at=$2 WHERE id=$3 RETURNING id`, status, timeNow, uint64(orderId)).Scan(&id)
	} else if status == cnst.Delivered {
		timeNow := time.Now().UTC().Format(cnst.Timestamptz)
		err = repo.db.QueryRowContext(ctx, `UPDATE "order" SET status=$1, delivered_at=$2 WHERE id=$3 RETURNING id`, status, timeNow, uint64(orderId)).Scan(&id)
	} else {
		err = repo.db.QueryRowContext(ctx, `UPDATE "order" SET status=$1 WHERE id=$2 RETURNING id`, status, uint64(orderId)).Scan(&id)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsOrderRelation
		}
		return 0, err
	}
	return alias.OrderId(id), nil
}

func (repo *RepoLayer) GetOrderSum(ctx context.Context, orderId alias.OrderId) (uint32, error) {
	var sum sql.NullInt32
	row := repo.db.QueryRowContext(ctx,
		`SELECT sum FROM "order" WHERE id=$1`, uint64(orderId))
	err := row.Scan(&sum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsOrderRelation
		}
		return 0, err
	}
	if !sum.Valid {
		return 0, myerrors.OrderSum
	}
	return uint32(sum.Int32), nil
}

func (repo *RepoLayer) GetFoodPrice(ctx context.Context, foodId alias.FoodId) (uint32, error) {
	var price uint32
	row := repo.db.QueryRowContext(ctx,
		`SELECT price FROM food WHERE id=$1`, uint64(foodId))
	err := row.Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsFoodRelation
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsFoodOrderRelation
		}
		return 0, err
	}
	return count, nil
}

func (repo *RepoLayer) UpdateSum(ctx context.Context, sum uint32, orderId alias.OrderId) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE "order" SET sum=$1 WHERE id=$2`, sum, uint64(orderId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.SqlNoRowsOrderRelation
	}
	return nil
}

func (repo *RepoLayer) AddToOrder(ctx context.Context, orderId alias.OrderId, foodId alias.FoodId, count uint32) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	res, err := repo.db.ExecContext(ctx,
		`INSERT INTO food_order (order_id, food_id, count,  created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`, uint64(orderId), uint64(foodId), count, timeNow, timeNow)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.OrderAddFood
	}
	sum, err := repo.GetOrderSum(ctx, orderId)
	if err != nil {
		if !errors.Is(err, myerrors.OrderSum) {
			return err
		}
		sum = 0
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
		return myerrors.SqlNoRowsFoodOrderRelation
	}
	price, err := repo.GetFoodPrice(ctx, foodId)
	if err != nil {
		return err
	}
	sum, err := repo.GetOrderSum(ctx, orderId)
	if err != nil {
		return err
	}
	if num := int(count) - int(currentCount); num > 0 {
		sum = sum + (count-currentCount)*price
	} else {
		sum = sum - (currentCount-count)*price
	}
	return repo.UpdateSum(ctx, sum, orderId)
}

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
		return myerrors.SqlNoRowsFoodOrderRelation
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
	if sum == 0 {
		res, err = repo.db.ExecContext(ctx,
			`DELETE FROM "order" WHERE id=$1`, uint64(orderId))
		if err != nil {
			return err
		}
		countRows, err = res.RowsAffected()
		if err != nil {
			return err
		}
		if countRows == 0 {
			return myerrors.SqlNoRowsOrderRelation
		}
		return nil
	}
	return repo.UpdateSum(ctx, sum, orderId)
}

func (repo *RepoLayer) CleanBasket(ctx context.Context, id alias.OrderId) error {
	res, err := repo.db.ExecContext(ctx,
		`DELETE FROM food_order WHERE order_id=$1`, uint64(id))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.FailCleanBasket
	}
	err = repo.UpdateSum(ctx, 0, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *RepoLayer) DeleteBasket(ctx context.Context, id alias.OrderId) error {
	res, err := repo.db.ExecContext(ctx,
		`DELETE FROM food_order WHERE order_id=$1`, uint64(id))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.FailCleanBasket
	}
	res, err = repo.db.ExecContext(ctx,
		`DELETE FROM "order" WHERE id=$1`, uint64(id))
	if err != nil {
		return err
	}
	countRows, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.FailCleanBasket
	}
	return nil
}

func (repo *RepoLayer) SetUser(ctx context.Context, orderId alias.OrderId, userId alias.UserId) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE "order" SET user_id=$1, unauth_token=NULL WHERE id=$2`, uint64(userId), uint64(orderId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.SqlNoRowsOrderRelation
	}
	return nil
}
