package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"2024_1_kayros/gen/go/comment"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	metrics "2024_1_kayros/microservices/metrics"
)

type Comment interface {
	Create(context.Context, *comment.Comment) (*comment.Comment, error)
	Delete(context.Context, *comment.CommentId) error
	GetCommentsByRest(context.Context, *comment.RestId) (*comment.CommentList, error)
}

type CommentLayer struct {
	db      *sql.DB
	metrics *metrics.MicroserviceMetrics
	stmt    map[string]*sql.Stmt
}

func NewCommentLayer(dbProps *sql.DB, metrics *metrics.MicroserviceMetrics, statements map[string]*sql.Stmt) Comment {
	return &CommentLayer{
		db:      dbProps,
		metrics: metrics,
		stmt: statements,
	}
}

func (repo *CommentLayer) Create(ctx context.Context, com *comment.Comment) (*comment.Comment, error) {
	var id uint64
	timeNow := time.Now()
	err := repo.stmt["createComment"].QueryRowContext(ctx, com.UserId, com.RestId, functions.MaybeNullString(com.Text), com.Rating).Scan(&id)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.INSERT).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsCommentRelation
		}
		return nil, err
	}
	com.Id = id

	rest := restCommentInfoDB{}
	timeNow = time.Now()
	row := repo.stmt["getRestRating"].QueryRowContext(ctx, com.RestId)
	timeEnd = time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	err = row.Scan(&rest.Rating, &rest.CommentCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsRestaurantRelation
		}
		return nil, err
	}
	r := fromNull(rest)
	newRating := (r.Rating*float64(r.CommentCount) + float64(com.Rating)) / (float64(r.CommentCount) + 1)
	timeNow = time.Now()
	res, err := repo.stmt["updateRestRating"].ExecContext(ctx, newRating, r.CommentCount+1, com.RestId)
	timeEnd = time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.UPDATE).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsRestaurantRelation
		}
		return nil, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if numRows == 0 {
		return nil, myerrors.SqlNoRowsUserRelation
	}
	timeNow = time.Now()
	res, err = repo.stmt["updateOrderCommentedStatus"].ExecContext(ctx, com.OrderId)
	timeEnd = time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.UPDATE).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsOrderRelation
		}
		return nil, err
	}
	numRows, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if numRows == 0 {
		return nil, myerrors.SqlNoRowsUserRelation
	}
	return com, err
}

func (repo *CommentLayer) GetCommentsByRest(ctx context.Context, restId *comment.RestId) (*comment.CommentList, error) {
	timeNow := time.Now()
	rows, err := repo.stmt["getRestComments"].QueryContext(ctx, restId.Id)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		return nil, err
	}
	comments := comment.CommentList{}
	for rows.Next() {
		com := comment.Comment{}
		err = rows.Scan(&com.Id, &com.UserName, &com.Image, &com.Text, &com.Rating)
		if err != nil {
			return nil, err
		}
		comments.Comment = append(comments.Comment, &com)
	}
	return &comments, nil
}

func (repo *CommentLayer) Delete(ctx context.Context, id *comment.CommentId) error {
	timeNow := time.Now()
	res, err := repo.stmt["deleteComment"].ExecContext(ctx, id.Id)
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.DELETE).Observe(float64(timeEnd.Milliseconds()))
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
	Rating       sql.NullFloat64
	CommentCount sql.NullInt32
}

type restCommentInfo struct {
	Rating       float64
	CommentCount uint32
}

func fromNull(r restCommentInfoDB) restCommentInfo {
	return restCommentInfo{
		Rating:       Float(r.Rating),
		CommentCount: Int(r.CommentCount),
	}
}

func Float(element sql.NullFloat64) float64 {
	fmt.Println(element)
	if element.Valid {
		return float64(element.Float64)
	}
	return 0
}

func Int(element sql.NullInt32) uint32 {
	fmt.Println(element)
	if element.Valid {
		return uint32(element.Int32)
	}
	return 0
}
