package food

import (
	"context"
	"database/sql"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Repo interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
	AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error)
	UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error)
	DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error)
}

type RepoLayer struct {
	DB *sql.DB
}

func NewRepoLayer(db *sql.DB) Repo {
	return &RepoLayer{DB: db}
}

func (repo *RepoLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	var food []*entity.Food
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, name, img_url, price, weight FROM food WHERE restaurant = $1", uint64(restId))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &entity.Food{}
		err = rows.Scan(item.Id, item.Name, item.ImgUrl, item.Price, item.Weight)
		if err != nil {
			return nil, err
		}
		food = append(food, item)
	}
	return food, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	item := &entity.Food{}
	row := repo.DB.QueryRowContext(ctx, "SELECT id, name, img_url, price, weight FROM food WHERE id=$1", foodId)
	err := row.Scan(item.Id, item.Name, item.ImgUrl, item.Price, item.Weight)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *RepoLayer) AddToOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	_, err := repo.DB.ExecContext(ctx, "INSERT INTO food-in-order (food_id, order_id, count) VALUES ($1, $2, 1)", foodId, orderId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *RepoLayer) UpdateCountInOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId, count uint32) (bool, error) {
	_, err := repo.DB.ExecContext(ctx, "UPDATE food-in-order SET count=$1 WHERE order_id=$2 AND food_id=$3", count, orderId, foodId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *RepoLayer) DeleteFromOrder(ctx context.Context, foodId alias.FoodId, orderId alias.OrderId) (bool, error) {
	_, err := repo.DB.ExecContext(ctx, "DELETE FROM food-in-order WHERE order_id=$1 AND food_id=$2", orderId, foodId)
	if err != nil {
		return false, err
	}
	return true, nil
}
