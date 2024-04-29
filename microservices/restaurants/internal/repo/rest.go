package repo

import (
	"context"
	"database/sql"
	"errors"

	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

type Repo interface {
	GetAll(ctx context.Context) ([]*rest.Rest, error)
	GetById(ctx context.Context, restId alias.RestId) (*rest.Rest, error)
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{
		db: dbProps,
	}
}

func (repo *RepoLayer) GetAll(ctx context.Context) ([]*rest.Rest, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, short_description, address, img_url FROM restaurant`)
	if err != nil {
		return nil, err
	}
	rests := []*rest.Rest{}
	for rows.Next() {
		Rest := rest.Rest{}
		err = rows.Scan(&Rest.Id, &Rest.Name, &Rest.ShortDescription, &Rest.Address, &Rest.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests = append(rests, &Rest)
	}
	return rests, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, restId alias.RestId) (*rest.Rest, error) {
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, long_description, address, img_url FROM restaurant WHERE id=$1`, uint64(restId))
	Rest := rest.Rest{}
	err := row.Scan(&Rest.Id, &Rest.Name, &Rest.LongDescription, &Rest.Address, &Rest.ImgUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsRestaurantRelation
		}
		return nil, err
	}
	return &Rest, nil
}
