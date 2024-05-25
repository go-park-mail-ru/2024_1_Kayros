package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/internal/utils/myerrors"
	metrics "2024_1_kayros/microservices/metrics"
)

type Rest interface {
	GetAll(ctx context.Context) (*rest.RestList, error)
	GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error)
	GetByFilter(ctx context.Context, id *rest.Id) (*rest.RestList, error)
	GetCategoryList(ctx context.Context) (*rest.CategoryList, error)
}

type RestLayer struct {
	db *sql.DB
	metrics *metrics.MicroserviceMetrics
}

func NewRestLayer(dbProps *sql.DB, metrics *metrics.MicroserviceMetrics) Rest {
	return &RestLayer{
		db: dbProps,
		metrics: metrics,
	}
}

func (repo *RestLayer) GetAll(ctx context.Context) (*rest.RestList, error) {
	timeNow := time.Now()
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, short_description, address, img_url FROM restaurant ORDER BY rating DESC`)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		return nil, err
	}
	rests := rest.RestList{}
	for rows.Next() {
		r := rest.Rest{}
		err = rows.Scan(&r.Id, &r.Name, &r.ShortDescription, &r.Address, &r.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests.Rest = append(rests.Rest, &r)
	}
	return &rests, nil
}


func (repo *RestLayer) GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error) {
	timeNow := time.Now()
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, long_description, address, img_url, rating, comment_count FROM restaurant WHERE id=$1`, id.Id)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	r := rest.Rest{}
	err := row.Scan(&r.Id, &r.Name, &r.LongDescription, &r.Address, &r.ImgUrl, &r.Rating, &r.CommentCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsRestaurantRelation
		}
		return nil, err
	}
	return &r, nil
}

func (repo *RestLayer) GetByFilter(ctx context.Context, id *rest.Id) (*rest.RestList, error) {
	timeNow := time.Now()
	rows, err := repo.db.QueryContext(ctx,
		`SELECT r.id, r.name, r.short_description, r.img_url FROM restaurant as r 
				JOIN rest_categories AS rc ON r.id=rc.restaurant_id WHERE rc.category_id=$1`, id.Id)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		return nil, err
	}
	rests := rest.RestList{}
	for rows.Next() {
		r := rest.Rest{}
		err = rows.Scan(&r.Id, &r.Name, &r.ShortDescription, &r.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests.Rest = append(rests.Rest, &r)
	}
	if len(rests.GetRest()) == 0 {
		return nil, nil
	}
	return &rests, nil
}

func (repo *RestLayer) GetCategoryList(ctx context.Context) (*rest.CategoryList, error) {
	timeNow := time.Now()
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name FROM category WHERE type='rest'`)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		return nil, err
	}
	categories := rest.CategoryList{}
	for rows.Next() {
		cat := rest.Category{}
		err = rows.Scan(&cat.Id, &cat.Name)
		if err != nil {
			return nil, err
		}
		categories.C = append(categories.C, &cat)
	}
	return &categories, nil
}
