package statistic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"2024_1_kayros/internal/delivery/metrics"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"

	"2024_1_kayros/internal/entity"
)

type Repo interface {
	Create(ctx context.Context, questionId uint64, rating uint8, userId string) error
	GetQuestionsOnFocus(ctx context.Context, url string) ([]*entity.Question, error)
	GetQuestions(ctx context.Context) ([]*entity.Question, error)
	NPS(ctx context.Context, id uint64) (int8, error)
	CSAT(ctx context.Context, id uint64) (int8, error)
}

type RepoLayer struct {
	db      *sql.DB
	metrics *metrics.Metrics
}

func NewRepoLayer(dbProps *sql.DB, metrics *metrics.Metrics) Repo {
	return &RepoLayer{
		db:      dbProps,
		metrics: metrics,
	}
}

func (repo *RepoLayer) Create(ctx context.Context, questionId uint64, rating uint8, user string) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	timeNowMetrics := time.Now()
	res, err := repo.db.ExecContext(ctx, `INSERT INTO quiz(question_id, user_id, rating, created_at) VALUES($1, $2, $3, $4)`, questionId, user, rating, timeNow)
	msRequestTimeout := time.Since(timeNowMetrics)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.INSERT).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.QuizAdd
	}
	return nil
}

func (repo *RepoLayer) GetQuestionsOnFocus(ctx context.Context, url string) ([]*entity.Question, error) {
	timeNow := time.Now()
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name, url, focus_id, param_type FROM question WHERE url=$1`, url)
	msRequestTimeout := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
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
	timeNow := time.Now()
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name, param_type FROM question`)
	msRequestTimeout := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
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
	var promoters uint64 // rating 9-10
	var critics uint64   // rating 0-6
	var respondents uint64

	timeNow := time.Now()
	row := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`, id)
	msRequestTimeout := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
	err := row.Scan(&promoters)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		promoters = 0
	}

	timeNow = time.Now()
	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating<7 AND question_id=$1`, id)
	msRequestTimeout = time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
	err = row.Scan(&critics)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		critics = 0
	}

	timeNow = time.Now()
	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE question_id=$1`, id)
	msRequestTimeout = time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
	err = row.Scan(&respondents)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		respondents = 0
	}

	if respondents == 0 {
		return 0, nil
	}
	return int8((promoters - critics) * 100 / respondents), nil
}

func (repo *RepoLayer) CSAT(ctx context.Context, id uint64) (int8, error) {
	var promoters uint64 // rating 9-10
	var respondents uint64

	timeNow := time.Now()
	row := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`, id)
	msRequestTimeout := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
	err := row.Scan(&promoters)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		promoters = 0
	}

	timeNow = time.Now()
	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE question_id=$1`, id)
	msRequestTimeout = time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(cnst.SELECT).Observe(float64(msRequestTimeout.Milliseconds()))
	err = row.Scan(&respondents)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		respondents = 0
	}

	if respondents == 0 {
		return 0, nil
	}
	return int8(promoters * 100 / respondents), nil
}
