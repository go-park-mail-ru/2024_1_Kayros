package food

import (
	"context"
	"database/sql"

	"2024_1_kayros/internal/entity"
)

type RepoInterface interface {
	GetByRest(ctx context.Context, restId int) ([]*entity.Food, error)
	GetById(ctx context.Context, id int) (*entity.Food, error)
	AddToOrder(ctx context.Context, foodId int, orderId int) error
	UpdateCountInOrder(ctx context.Context, foodId int, orderId int, count int) error
	DeleteFromOrder(ctx context.Context, foodId int, orderId int) error
}

type Repo struct {
	DB *sql.DB
}

func NewFoodRepository(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (repo *Repo) GetByRest(ctx context.Context, restId int) ([]*entity.Food, error) {
	var food []*entity.Food
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, name, img_url, price, weight FROM food WHERE restaurant = $1", restId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := &entity.Food{}
		err = rows.Scan(&item.Id, &item.Name, &item.ImgUrl, &item.Price, &item.Weight)
		if err != nil {
			return nil, err
		}
		food = append(food, item)
	}
	return food, nil
}

func (repo *Repo) GetById(ctx context.Context, id int) (*entity.Food, error) {
	item := &entity.Food{}
	row := repo.DB.QueryRowContext(ctx, "SELECT id, name, img_url, price, weight FROM food WHERE id=$1", id)
	err := row.Scan(&item.Id, &item.Name, &item.ImgUrl, &item.Price, &item.Weight)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *Repo) AddToOrder(ctx context.Context, foodId int, orderId int) error {
	_, err := repo.DB.ExecContext(ctx, "INSERT INTO food-in-order (food_id, order_id, count) VALUES ($1, $2, 1)", foodId, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repo) UpdateCountInOrder(ctx context.Context, foodId int, orderId int, count int) error {
	_, err := repo.DB.ExecContext(ctx, "UPDATE food-in-order SET count=$1 WHERE order_id=$2 AND food_id=$3", count, orderId, foodId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repo) DeleteFromOrder(ctx context.Context, foodId int, orderId int) error {
	_, err := repo.DB.ExecContext(ctx, "DELETE FROM food-in-order WHERE order_id=$1 AND food_id=$2", orderId, foodId)
	if err != nil {
		return err
	}
	return nil
}
