package search

import (
	"context"
	"database/sql"
	"errors"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/utils/myerrors"
)

type Repo interface {
	Search(ctx context.Context, search string) ([]*dto.RestaurantAndFood, error)
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(db *sql.DB) Repo {
	return &RepoLayer{
		db: db,
	}
}

func (repo *RepoLayer) Search(ctx context.Context, search string) ([]*dto.RestaurantAndFood, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, img_url FROM restaurant 
			WHERE LOWER(name) LIKE LOWER('%' || $1 || '%')`, search)
	if err != nil {
		return nil, err
	}
	rests := []*dto.RestaurantAndFood{}
	rests, err = repo.SelectRests(ctx, rows, rests)
	if err != nil {
		return nil, err
	}
	if len(rests) == 0 {
		rows, err = repo.db.QueryContext(ctx,
			`SELECT DISTINCT r.id, r.name, r.img_url FROM restaurant AS r
			JOIN rest_categories AS rc ON r.id=rc.restaurant_id JOIN category AS c
            ON rc.category_id=c.id WHERE LOWER(c.name) LIKE LOWER('%' || $1 || '%')`, search)
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
		rs, err := repo.db.QueryContext(ctx, `SELECT id, name FROM category AS c
			JOIN rest_categories AS rc ON c.id=rc.category_id WHERE rc.restaurant_id=$1`, rest.Id)
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
