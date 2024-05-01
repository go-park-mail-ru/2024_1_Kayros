package comment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
)

type Repo interface {
	Create(ctx context.Context, comment entity.Comment) (*entity.Comment, error)
	GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error)
	Delete(ctx context.Context, id uint64) error
}

type RepoLayer struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{
		db: dbProps,
	}
}

func (repo *RepoLayer) Create(ctx context.Context, comment entity.Comment) (*entity.Comment, error) {
	row := repo.db.QueryRowContext(ctx,
		`INSERT INTO "comment" (user_id, restaurant_id, text, rating) VALUES ($1, $2, $3, $4) RETURNING id, user_id, text, rating`, comment.UserId, comment.RestId, comment.Text, comment.Rating)
	res := entity.Comment{}
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

func (repo *RepoLayer) GetCommentsByRest(ctx context.Context, restId alias.RestId) ([]*entity.Comment, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT c.id, u.name, u.img_url, c.text, c.rating FROM "comment" AS c JOIN "user" AS u ON c.user_id = u.id WHERE restaurant_id=$1`, restId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	comments := []*entity.Comment{}
	for rows.Next() {
		com := entity.Comment{}
		err = rows.Scan(&com.Id, &com.UserName, &com.UserImage, &com.Text, &com.Rating)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		comments = append(comments, &com)
	}
	return comments, nil
}

func (repo *RepoLayer) Delete(ctx context.Context, id uint64) error {
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
