package statistic

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
)

type Repo interface {
	Create(ctx context.Context, questionId uint64, rating uint32, userId string, time string) error
	Update(ctx context.Context, questionId uint64, rating uint32, userId string) error
	//GetStatistic(ctx context.Context) ([]*entity.Statistic, error)
	GetQuestionInfo(ctx context.Context) ([]*entity.Question, error)
	NPS(ctx context.Context, id uint16) (int8, error)
	CSAT(ctx context.Context, id uint16) (int8, error)
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
	res, err := repo.db.ExecContext(ctx,
		`UPDATE quiz SET rating=$1 WHERE question_id=$2 AND user_id=$3`, rating, questionId, user)
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
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name, url, focus_id, param_type FROM question`)
	if err != nil {
		return nil, err
	}
	qs := []*entity.Question{}
	for rows.Next() {
		q := entity.Question{}
		err = rows.Scan(&q.Id, &q.Name, &q.Url, &q.FocusId, &q.ParamType)
		if err != nil {
			return nil, err
		}
		qs = append(qs, &q)
	}
	return qs, nil
}

//func (repo *RepoLayer) GetStatistic(ctx context.Context) ([]*entity.Statistic, error) {
//	rows, err := repo.db.QueryContext(ctx, `SELECT question_id, COUNT(*), AVG(rating) FROM quiz JOIN question ON quiz.question_id = question.id  GROUP BY question_id`)
//	if err != nil {
//		return nil, err
//	}
//	stats := []*entity.Statistic{}
//	for rows.Next() {
//		stat := entity.Statistic{}
//		err = rows.Scan(&stat.QuestionId, &stat.Count, &stat.Rating, &stat.QuestionName)
//		if err != nil {
//			return nil, err
//		}
//		stats = append(stats, &stat)
//	}
//	return stats, nil
//}

func (repo *RepoLayer) NPS(ctx context.Context, id uint16) (int8, error) {
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

	return (first - second) * 100 / n, nil
}

func (repo *RepoLayer) CSAT(ctx context.Context, id uint16) (int8, error) {
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

	return top * 100 / n, nil
}
