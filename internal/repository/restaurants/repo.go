package restaurants

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"go.uber.org/zap"
)

const NoRestError = "Такого ресторана нет"

type Repo interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
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

func (repo *RepoLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, short_description, long_description, address, img_url FROM "Restaurant"`)
	if err != nil {
		return nil, err
	}

	var rests []*entity.Restaurant
	for rows.Next() {
		rest := entity.Restaurant{}
		err = rows.Scan(&rest.Id, &rest.Name, &rest.ShortDescription, &rest.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests = append(rests, &rest)
	}
	return rests, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, short_description, long_description, address, img_url FROM "Restaurant" WHERE id=$1`, uint64(restId))
	var rest *entity.Restaurant
	err := row.Scan(&rest.Id, &rest.Name, &rest.LongDescription, &rest.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf(NoRestError)
	}
	if err != nil {
		return nil, err
	}
	return rest, nil
}
