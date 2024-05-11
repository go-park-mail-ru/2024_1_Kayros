package statistic

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	db *sql.DB
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{
		db: dbProps,
	}
}

func (repo *RepoLayer) Create(ctx context.Context, questionId uint64, rating uint8, user string) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	res, err := repo.db.ExecContext(ctx, `INSERT INTO quiz(question_id, user_id, rating, created_at) VALUES($1, $2, $3, $4)`, questionId, user, rating, timeNow)
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
	var promoters uint64 // rating 9-10
	var critics uint64   // rating 0-6
	var respondents uint64

	row := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`, id)
	err := row.Scan(&promoters)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		promoters = 0
	}

	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating<7 AND question_id=$1`, id)
	err = row.Scan(&critics)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		critics = 0
	}

	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE question_id=$1`, id)
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

	row := repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`, id)
	err := row.Scan(&promoters)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
		promoters = 0
	}

	row = repo.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz WHERE question_id=$1`, id)
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
