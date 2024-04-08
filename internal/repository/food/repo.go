package food

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

const NoFoodError = "Такого блюда нет"

type Repo interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
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

func (repo *RepoLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT f.id, f.name, description, restaurant_id, ca.name, weight, price, img_url FROM food as f 
    JOIN category as ca ON f.category_id=ca.id WHERE restaurant_id = $1`, uint64(restId))
	if err != nil {
		return nil, err
	}

	var food []*entity.Food
	for rows.Next() {
		item := entity.Food{}
		err = rows.Scan(&item.Id, &item.Name, &item.Description, &item.RestaurantId,
			&item.Category, &item.Weight, &item.Price, &item.ImgUrl)
		if err != nil {
			return nil, err
		}
		food = append(food, &item)
	}
	return food, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, description, restaurant_id, category_id, weight, price, img_url 
				FROM food WHERE id=$1`, uint64(foodId))

	var item *entity.Food
	err := row.Scan(&item.Id, &item.Name, &item.Description, &item.RestaurantId,
		&item.Category, &item.ImgUrl, &item.Price, &item.Weight)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf(NoFoodError)
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}
