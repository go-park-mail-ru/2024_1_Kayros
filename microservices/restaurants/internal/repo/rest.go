package repo

import (
	"context"
	"database/sql"
	"errors"

	"2024_1_kayros/internal/utils/myerrors"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

type Rest interface {
	GetAll(ctx context.Context) (*rest.RestList, error)
	GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error)
	GetByFilter(ctx context.Context, filter *rest.Filter) (*rest.RestList, error)
}

type RestLayer struct {
	db *sql.DB
}

func NewRestLayer(dbProps *sql.DB) Rest {
	return &RestLayer{
		db: dbProps,
	}
}

func (repo *RestLayer) GetAll(ctx context.Context) (*rest.RestList, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, short_description, address, img_url FROM restaurant`)
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
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, long_description, address, img_url FROM restaurant WHERE id=$1`, id)
	r := rest.Rest{}
	err := row.Scan(&r.Id, &r.Name, &r.LongDescription, &r.Address, &r.ImgUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsRestaurantRelation
		}
		return nil, err
	}
	return &r, nil
}

func (repo *RestLayer) GetByFilter(ctx context.Context, filter *rest.Filter) (*rest.RestList, error) {
	var id uint64
	err := repo.db.QueryRowContext(ctx, `SELECT id FROM category WHERE name=$1`, filter.Filter).Scan(&id)
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx,
		`SELECT r.id, r.name, r.short_description, r.img_url FROM restaurant as r 
				JOIN rest_categories AS rc ON r.id=rc.restaurant_id WHERE rc.category_id=$1`, id)
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
