package restaurants

import (
	"context"
	"database/sql"
	"errors"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
)

type Repo interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{
		db: dbProps,
	}
}

func (repo *RepoLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, short_description, address, img_url FROM restaurant`)
	if err != nil {
		return nil, err
	}
	rests := []*entity.Restaurant{}
	for rows.Next() {
		rest := entity.Restaurant{}
		err = rows.Scan(&rest.Id, &rest.Name, &rest.ShortDescription, &rest.Address, &rest.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests = append(rests, &rest)
	}
	return rests, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, long_description, address, img_url FROM restaurant WHERE id=$1`, uint64(restId))
	rest := entity.Restaurant{}
	err := row.Scan(&rest.Id, &rest.Name, &rest.LongDescription, &rest.Address, &rest.ImgUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsRestaurantRelation
		}
		return nil, err
	}
	return &rest, nil
}
