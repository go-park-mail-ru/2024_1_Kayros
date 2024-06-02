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

//go:generate mockgen -source ./rest.go -destination=./mocks/service.go -package=mock_service
type Rest interface {
	GetAll(ctx context.Context) (*rest.RestList, error)
	GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error)
	GetByFilter(ctx context.Context, id *rest.Id) (*rest.RestList, error)
	GetCategoryList(ctx context.Context) (*rest.CategoryList, error)
	GetTop(ctx context.Context, limit uint64) (*rest.RestList, error)
	GetLastRests(ctx context.Context, userId uint64, limit uint64) (*rest.RestList, error)
}

type RestLayer struct {
	db      *sql.DB
	metrics *metrics.MicroserviceMetrics
	stmt    map[string]*sql.Stmt
}

func NewRestLayer(dbProps *sql.DB, metrics *metrics.MicroserviceMetrics, statements map[string]*sql.Stmt) Rest {
	return &RestLayer{
		db:      dbProps,
		metrics: metrics,
		stmt: statements,
	}
}

func (repo *RestLayer) GetAll(ctx context.Context) (*rest.RestList, error) {
	timeNow := time.Now()
	rows, err := repo.stmt["getAll"].QueryContext(ctx)
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
	row := repo.stmt["getRestById"].QueryRowContext(ctx, id.Id)
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
	rows, err := repo.stmt["getRestListUsingFilter"].QueryContext(ctx, id.Id)
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
	rows, err := repo.stmt["getRestsByCategory"].QueryContext(ctx)
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

func (repo *RestLayer) GetTop(ctx context.Context, limit uint64) (*rest.RestList, error) {
	rows, err := repo.stmt["getTopRests"].QueryContext(ctx, limit)
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
	return &rests, nil
}

func (repo *RestLayer) GetLastRests(ctx context.Context, userId uint64, limit uint64) (*rest.RestList, error) {
	rows, err := repo.stmt["getLastRests"].QueryContext(ctx, userId, limit)
	if err != nil {
		return nil, err
	}
	//получили список id ресторанов, из которых в последние разы заказывал человек
	rests := rest.RestList{}
	for rows.Next() {
		r := rest.Id{}
		err = rows.Scan(&r.Id)
		if err != nil {
			return nil, err
		}
		rest := rest.Rest{}
		//для каждого получили инфу
		err = repo.stmt["getShortRestById"].QueryRowContext(ctx, r.Id).Scan(&rest.Id, &rest.Name, &rest.ShortDescription, &rest.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests.Rest = append(rests.Rest, &rest)
	}
	return &rests, nil
}
