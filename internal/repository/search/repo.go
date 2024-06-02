package search

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity/dto"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
)

type Repo interface {
	Search(ctx context.Context, search string) ([]*dto.RestaurantAndFood, error)
}

type RepoLayer struct {
	db *sql.DB
	metrics *metrics.Metrics
	stmt map[string]*sql.Stmt
}

func NewRepoLayer(db *sql.DB, metrics *metrics.Metrics, statements map[string]*sql.Stmt) Repo {
	return &RepoLayer{
		db: db,
		metrics: metrics,
		stmt: statements,
	}
}

func (repo *RepoLayer) Search(ctx context.Context, search string) ([]*dto.RestaurantAndFood, error) {
	timeNow := time.Now()
	rows, err := repo.stmt["selectRestBySearch"].QueryContext(ctx, search)
	msRequestTimeout := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		return nil, err
	}
	rests := []*dto.RestaurantAndFood{}
	rests, err = repo.SelectRests(ctx, rows, rests)
	if err != nil {
		return nil, err
	}
	if len(rests) == 0 {
		timeNow = time.Now()
		rows, err = repo.stmt["getRestsByCategory"].QueryContext(ctx, search)
		msRequestTimeout = time.Since(timeNow)
		repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
		if err != nil {
			return nil, err
		}
		rests, err = repo.SelectRests(ctx, rows, rests)
		if err != nil {
			return nil, err
		}
	}
	return rests, nil
}

func (repo *RepoLayer) SelectRests(ctx context.Context, rows *sql.Rows, rests []*dto.RestaurantAndFood) ([]*dto.RestaurantAndFood, error) {
	var err error
	for rows.Next() {
		var rest dto.RestaurantAndFood
		err = rows.Scan(&rest.Id, &rest.Name, &rest.ImgUrl)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			return nil, err
		}
		timeNow := time.Now()
		rs, err := repo.stmt["selectRests"].QueryContext(ctx, rest.Id)
		msRequestTimeout := time.Since(timeNow)
		repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, myerrors.SqlNoRowsUserRelation
			}
			return nil, err
		}
		for rs.Next() {
			var cat dto.Category
			err = rs.Scan(&cat.Id, &cat.Name)
			if err != nil {
				return nil, err
			}
			rest.Categories = append(rest.Categories, &cat)
		}
		rests = append(rests, &rest)
	}
	return rests, nil
}
