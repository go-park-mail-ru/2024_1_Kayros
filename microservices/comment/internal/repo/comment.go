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
	var id uint64
	err := repo.db.QueryRowContext(ctx,
		`INSERT INTO "comment" (user_id, restaurant_id, text, rating) VALUES ($1, $2, $3, $4) RETURNING id`, com.UserId, com.RestId, com.Text, com.Rating).Scan(&id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsCommentRelation
		}
		return nil, err
	}
	com.Id = id
	//rest := restCommentInfoDB{}
	//err = repo.db.QueryRowContext(ctx,
	//	`SELECT rating, comment_count FROM restaurant WHERE id=$1`, com.RestId).Scan(&rest.Rating, &rest.CommentCount)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return nil, myerrors.SqlNoRowsCommentRelation
	//	}
	//	return nil, err
	//}
	//r := fromNull(rest)
	//newRating := (r.Rating*r.CommentCount + com.Rating) / (r.CommentCount + 1)
	//row, err := repo.db.ExecContext(ctx,
	//	`UPDATE restaurant SET rating=$1, comment_count=$2 WHERE id=$3`, newRating, r.CommentCount+1, com.RestId)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return nil, myerrors.SqlNoRowsCommentRelation
	//	}
	//	return nil, err
	//}
	//numRows, err := row.RowsAffected()
	//if err != nil {
	//	return nil, err
	//}
	//if numRows == 0 {
	//	return nil, myerrors.SqlNoRowsUserRelation
	//}
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return nil, myerrors.SqlNoRowsCommentRelation
	//	}
	//	return nil, err
	//}
	return com, err
}

func (repo *CommentLayer) GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT c.id, u.name, u.img_url, c.text, c.rating FROM "comment" AS c JOIN "user" AS u ON c.user_id = u.id WHERE restaurant_id=$1`, restId.Id)
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
		`DELETE FROM "comment" WHERE id=$1`, id.Id)
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

type restCommentInfoDB struct {
	Rating       sql.NullInt32
	CommentCount sql.NullInt32
}

type restCommentInfo struct {
	Rating       uint32
	CommentCount uint32
}

func fromNull(r restCommentInfoDB) restCommentInfo {
	return restCommentInfo{
		Rating:       Int(r.Rating),
		CommentCount: Int(r.CommentCount),
	}
}

func Int(element sql.NullInt32) uint32 {
	fmt.Println(element)
	if element.Valid {
		return uint32(element.Int32)
	}
	return 0
}
