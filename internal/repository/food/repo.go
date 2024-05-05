package food

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"2024_1_kayros/internal/utils/myerrors"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Repo interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
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

func (repo *RepoLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT c.name, f.id, f.name, restaurant_id, weight, price, img_url FROM food as f
   JOIN category as c ON f.category_id=c.id WHERE restaurant_id = $1 ORDER BY category_id`, uint64(restId))
	if err != nil {
		return nil, err
	}
	food := []*entity.Food{}
	for rows.Next() {
		item := entity.Food{}
		err = rows.Scan(&item.Category, &item.Id, &item.Name, &item.RestaurantId,
			&item.Weight, &item.Price, &item.ImgUrl)
		if err != nil {
			return nil, err
		}
		food = append(food, &item)
	}
	return food, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, restaurant_id, category_id, weight, price, img_url
				FROM food WHERE id=$1`, uint64(foodId))
	var item entity.Food
	err := row.Scan(&item.Id, &item.Name, &item.RestaurantId,
		&item.Category, &item.Weight, &item.Price, &item.ImgUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsFoodRelation
		}
		return nil, err
	}
	return &item, nil
}
