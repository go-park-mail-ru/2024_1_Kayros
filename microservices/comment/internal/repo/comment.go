package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"2024_1_kayros/internal/utils/myerrors"
	comment "2024_1_kayros/microservices/comment/proto"
)

type Comment interface {
	Create(context.Context, *comment.Comment) (*comment.Comment, error)
	Delete(context.Context, *comment.CommentId) error
	GetCommentsByRest(context.Context, *comment.RestId) (*comment.CommentList, error)
}

type CommentLayer struct {
	db *sql.DB
}

func NewCommentLayer(dbProps *sql.DB) Comment {
	return &CommentLayer{
		db: dbProps,
	}
}

func (repo *CommentLayer) Create(ctx context.Context, com *comment.Comment) (*comment.Comment, error) {
	row := repo.db.QueryRowContext(ctx,
		`INSERT INTO "comment" (user_id, restaurant_id, text, rating) VALUES ($1, $2, $3, $4) RETURNING id, user_id, text, rating`, com.UserId, com.RestId, com.Text, com.Rating)
	res := comment.Comment{}
	err := row.Scan(&res.Id, &res.UserId, &res.RestId, &res.Text, &res.Rating)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsCommentRelation
		}
		return nil, err
	}
	return &res, err
}

func (repo *CommentLayer) GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT c.id, u.name, u.img_url, c.text, c.rating FROM "comment" AS c JOIN "user" AS u ON c.user_id = u.id WHERE restaurant_id=$1`, restId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	comments := comment.CommentList{}
	for rows.Next() {
		com := comment.Comment{}
		err = rows.Scan(&com.Id, &com.UserName, &com.Image, &com.Text, &com.Rating)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		comments.Comment = append(comments.Comment, &com)
	}
	return &comments, nil
}

func (repo *CommentLayer) Delete(ctx context.Context, id *comment.CommentId) error {
	res, err := repo.db.ExecContext(ctx,
		`DELETE FROM "comment" WHERE id=$1`, id)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.SqlNoRowsCommentRelation
	}
	return nil
}
