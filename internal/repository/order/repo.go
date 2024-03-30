package order

import (
	"context"
	"database/sql"

	"2024_1_kayros/internal/entity"
)

type RepoInterface interface {
	Create(ctx context.Context, user uint64, date string, status string) error
	GetBasket(ctx context.Context, userId uint64, status string) (*entity.Order, error)
	GetBasketId(ctx context.Context, userId uint64, status string) (uint64, error)
	GetFood(ctx context.Context, id uint64) ([]*entity.FoodInOrder, error)
	Update(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, orderId uint64, status string) (string, error)
	AddToOrder(ctx context.Context, orderId uint64, foodId int, count int) error
	UpdateCountInOrder(ctx context.Context, orderId uint64, foodId int, count int) error
	DeleteFromOrder(ctx context.Context, orderId uint64, foodId int) error
}

type Repo struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) RepoInterface {
	return &Repo{DB: db}
}

func (repo *Repo) Create(ctx context.Context, user uint64, date string, status string) error {
	_, err := repo.DB.ExecContext(ctx, "INSERT INTO order(user, date_order, status) VALUES ($1, $2, $3)", user, date, status)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repo) GetBasket(ctx context.Context, userId uint64, status string) (*entity.Order, error) {
	row := repo.DB.QueryRowContext(ctx, "SELECT * FROM order WHERE user_id= $1 AND status=$2", userId, status)
	var order *entity.Order
	err := row.Scan(&order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *Repo) GetBasketId(ctx context.Context, userId uint64, status string) (uint64, error) {
	row := repo.DB.QueryRowContext(ctx, "SELECT id FROM order WHERE user_id= $1 AND status=$2", userId, status)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *Repo) GetFood(ctx context.Context, id uint64) ([]*entity.FoodInOrder, error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT food.food_id, name, img_url, weight, price, count FROM food_order INNER JOIN food ON food_order.food_id==food.id WHERE order_id= $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var foodArray []*entity.FoodInOrder
	for rows.Next() {
		var food *entity.FoodInOrder
		err = rows.Scan(&food.Id, &food.Name, &food.ImgUrl, &food.Weight, &food.Price, &food.Count)
		if err != nil {
			return nil, err
		}
		foodArray = append(foodArray, food)
	}
	return foodArray, nil
}

func (repo *Repo) Update(ctx context.Context, order *entity.Order) error {
	_, err := repo.DB.ExecContext(ctx, "UPDATE order SET date_receiving=$1, address=$2, extra_address=$3 WHERE order_id=$4", order.DateReceiving, order.Address, order.ExtraAddress)
	return err
}

func (repo *Repo) UpdateStatus(ctx context.Context, orderId uint64, status string) (string, error) {
	_, err := repo.DB.ExecContext(ctx, "UPDATE order SET status=$1 WHERE order_id=$2", status, orderId)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (repo *Repo) AddToOrder(ctx context.Context, orderId uint64, foodId int, count int) error {
	_, err := repo.DB.ExecContext(ctx, "INSERT INTO food_order (order_id, food_id, count) VALUES ($1, $2, $3)", orderId, foodId, count)
	return err
}

func (repo *Repo) UpdateCountInOrder(ctx context.Context, orderId uint64, foodId int, count int) error {
	_, err := repo.DB.ExecContext(ctx, "UPDATE food_order SET count=$1 WHERE order_id=$2 AND food_id=$3", count, orderId, foodId)
	return err
}

func (repo *Repo) DeleteFromOrder(ctx context.Context, orderId uint64, foodId int) error {
	_, err := repo.DB.ExecContext(ctx, "DELETE FROM food_order WHERE order_id=$1 AND food_id=$2", orderId, foodId)
	return err
}
