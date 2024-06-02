package food

import (
	"context"
	"database/sql"
	"errors"

	"2024_1_kayros/internal/utils/myerrors"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Repo interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
}

type Layer struct {
	db     *sql.DB
	stmt    map[string]*sql.Stmt
}

func NewLayer(dbProps *sql.DB, statements map[string]*sql.Stmt) Repo {
	return &Layer{
		db: dbProps,
		stmt: statements,
	}
}

func (repo *Layer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	rows, err := repo.stmt["getByRestId"].QueryContext(ctx, uint64(restId))
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

func (repo *Layer) GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error) {
	row := repo.stmt["getById"].QueryRowContext(ctx, uint64(foodId))
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
