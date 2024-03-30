package restaurants

import (
	"context"
	"database/sql"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Repo interface {
	GetAll(context.Context) ([]*entity.Restaurant, error)
	GetById(context.Context, alias.RestId) (*entity.Restaurant, error)
}

type RepoLayer struct {
	DB *sql.DB
}

func NewRepoLayer(db *sql.DB) Repo {
	return &RepoLayer{DB: db}
}

func (repo *RepoLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	var rests []*entity.Restaurant
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, name, short_description, img_url FROM Restaurant")
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rest := &entity.Restaurant{}
		err = rows.Scan(rest.Id, rest.Name, rest.ShortDescription, rest.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests = append(rests, rest)
	}
	return rests, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, id alias.RestId) (*entity.Restaurant, error) {
	rest := &entity.Restaurant{}
	row := repo.DB.QueryRowContext(ctx, "SELECT id, name, long_description, img_url FROM Restaurant WHERE id=$1", uint64(id))
	err := row.Scan(rest.Id, rest.Name, rest.LongDescription, rest.ImgUrl)
	if err != nil {
		return nil, err
	}
	return rest, nil
}
