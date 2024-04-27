package statistic

import (
	"context"
	"database/sql"
	"time"

	"2024_1_kayros/internal/entity"
	"go.uber.org/zap"
)

type Repo interface {
	Create(ctx context.Context, questionId uint64, rating uint32, userId string) error
	Update(ctx context.Context, questionId uint64, rating uint32, userId string) error
	GetStatistic(ctx context.Context) ([]*entity.Statistic, error)
	GetQuestionInfo(ctx context.Context) ([]*entity.Question, error)
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

func (repo *RepoLayer) Create(ctx context.Context, questionId uint64, rating uint32, user string) error {
	timeNow := time.Now().UTC().Format("2006-01-02 15:04:05-07:00")
	res, err := repo.db.ExecContext(ctx, `INSERT INTO quiz(question_id, user_id, rating, created_at) VALUES($1, $2, $3, $4)`, questionId, user, rating, timeNow)
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

func (repo *RepoLayer) GetQuestionInfo(ctx context.Context) ([]*entity.Question, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name FROM question`)
	if err != nil {
		return nil, err
	}
	qs := []*entity.Question{}
	for rows.Next() {
		q := entity.Question{}
		err = rows.Scan(&q.Id, &q.Name)
		if err != nil {
			return nil, err
		}
		qs = append(qs, &q)
	}
	return qs, nil
}

func (repo *RepoLayer) GetStatistic(ctx context.Context) ([]*entity.Statistic, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT question_id, COUNT(*), AVG(rating), MIN(name) as name FROM quiz JOIN question ON quiz.question_id = question.id  GROUP BY question_id;`)
	if err != nil {
		return nil, err
	}
	stats := []*entity.Statistic{}
	for rows.Next() {
		stat := entity.Statistic{}
		err = rows.Scan(&stat.QuestionId, &stat.Count, &stat.AverageRating, &stat.QuestionName)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &stat)
	}
	return stats, nil
}