package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"2024_1_kayros/internal/utils/myerrors"
	metrics "2024_1_kayros/microservices/metrics"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Repo interface {
	GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, foodId alias.FoodId) (*entity.Food, error)
}

type Layer struct {
	db      *sql.DB
	metrics *metrics.MicroserviceMetrics
}

func NewLayer(dbProps *sql.DB, metrics *metrics.MicroserviceMetrics) Repo {
	return &Layer{
		db:      dbProps,
		metrics: metrics,
	}
}

func (repo *Layer) GetByRestId(ctx context.Context, restId alias.RestId) ([]*entity.Food, error) {
	timeNow := time.Now()
	rows, err := repo.db.QueryContext(ctx,
		`SELECT c.name, f.id, f.name, restaurant_id, weight, price, img_url FROM food as f
   JOIN category as c ON f.category_id=c.id WHERE restaurant_id = $1 ORDER BY category_id`, uint64(restId))
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
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
	timeNow := time.Now()
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, restaurant_id, category_id, weight, price, img_url
				FROM food WHERE id=$1`, uint64(foodId))
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
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
