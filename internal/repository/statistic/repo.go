package statistic

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
)

type Repo interface {
	Create(ctx context.Context, questionId uint64, rating uint32, userId string, time string) error
	Update(ctx context.Context, questionId uint64, rating uint32, userId string) error
	GetQuestionInfo(ctx context.Context, url string) ([]*entity.Question, error)
	GetQuestions(ctx context.Context) ([]*entity.Question, error)
	NPS(ctx context.Context, id uint64) (int8, error)
	CSAT(ctx context.Context, id uint64) (int8, error)
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

func (repo *RepoLayer) Create(ctx context.Context, questionId uint64, rating uint32, user string, time string) error {
	res, err := repo.db.ExecContext(ctx, `INSERT INTO quiz(question_id, user_id, rating, created_at) VALUES($1, $2, $3, $4)`, questionId, user, rating, time)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return err
	}
	return nil
}

func (repo *RepoLayer) Update(ctx context.Context, questionId uint64, rating uint32, user string) error {
	timeNow := time.Now().UTC().Format("2006-01-02 15:04:05-07:00")
	res, err := repo.db.ExecContext(ctx,
		`UPDATE quiz SET rating=$1, created_at=$2 WHERE question_id=$3 AND user_id=$4`, rating, timeNow, questionId, user)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return err
	}
	return nil
}

func (repo *RepoLayer) GetQuestionInfo(ctx context.Context, url string) ([]*entity.Question, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name, url, focus_id, param_type FROM question WHERE url=$1`, url)
	if err != nil {
		return nil, err
	}
	qs := []*entity.Question{}
	for rows.Next() {
		qSql := entity.QuestionSql{}
		err = rows.Scan(&qSql.Id, &qSql.Name, &qSql.Url, &qSql.FocusId, &qSql.ParamType)
		if err != nil {
			return nil, err
		}
		q := entity.QuestionFromDB(&qSql)
		qs = append(qs, q)
	}
	return qs, nil
}

func (repo *RepoLayer) GetQuestions(ctx context.Context) ([]*entity.Question, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name, param_type FROM question`)
	if err != nil {
		return nil, err
	}
	qs := []*entity.Question{}
	for rows.Next() {
		q := entity.Question{}
		err = rows.Scan(&q.Id, &q.Name, &q.ParamType)
		if err != nil {
			return nil, err
		}
		qs = append(qs, &q)
	}
	return qs, nil
}

func (repo *RepoLayer) NPS(ctx context.Context, id uint64) (int8, error) {
	//кол-во промоутеров 9-10
	var first int8
	var second int8
	var n int8

	row := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`, id)
	err := row.Scan(&first)
	if err != nil {
		return 0, err
	}

	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating<7 AND question_id=$1`, id)
	err = row.Scan(&second)
	if err != nil {
		return 0, err
	}
	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE question_id=$1`, id)
	err = row.Scan(&n)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}
	return (first - second) * 100 / n, nil
}

func (repo *RepoLayer) CSAT(ctx context.Context, id uint64) (int8, error) {
	//кол-во промоутеров 9-10
	var top int8
	var n int8

	row := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating>9 AND question_id=$1`, id)
	err := row.Scan(&top)
	if err != nil {
		return 0, err
	}

	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE question_id=$1`, id)
	err = row.Scan(&n)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}
	return top * 100 / n, nil
}
