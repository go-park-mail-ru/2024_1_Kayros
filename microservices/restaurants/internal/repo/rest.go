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
	CreateComment(context.Context, *rest.Comment) (*rest.Comment, error)
	DeleteComment(context.Context, *rest.CommentId) (*rest.Empty, error)
	GetCommentsByRest(context.Context, *rest.RestId) (*rest.CommentList, error)
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

func (repo *RestLayer) CreateComment(context.Context, *rest.Comment) (*rest.Comment, error) {
	return nil, nil
}
func (repo *RestLayer) DeleteComment(context.Context, *rest.CommentId) (*rest.Empty, error) {
	return nil, nil
}
func (repo *RestLayer) GetCommentsByRest(context.Context, *rest.RestId) (*rest.CommentList, error) {
	return nil, nil
}
