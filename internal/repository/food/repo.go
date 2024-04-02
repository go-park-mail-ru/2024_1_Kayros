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
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{db: dbProps}
}

func (repo *RepoLayer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	var food []*entity.Food
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, img_url, price, weight FROM food WHERE restaurant = $1", uint64(restId))
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
	row := repo.db.QueryRowContext(ctx, "SELECT id, name, img_url, price, weight FROM food WHERE id=$1", foodId)
	err := row.Scan(item.Id, item.Name, item.ImgUrl, item.Price, item.Weight)
	if err != nil {
		return nil, err
	}
	return item, nil
}
