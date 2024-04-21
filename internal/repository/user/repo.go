package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"

	"2024_1_kayros/internal/entity"
)

type Repo interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	DeleteByEmail(ctx context.Context, email string) error

	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, uDataChange *entity.User, email string) error
}

type RepoLayer struct {
	database *sql.DB
}

func NewRepoLayer(db *sql.DB) Repo {
	return &RepoLayer{
		database: db,
	}
}

func (repo *RepoLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, email, COALESCE(phone, ''), password, COALESCE(address, ''), img_url, COALESCE(card_number, '')  FROM "user" WHERE email = $1`, email)
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Address, &user.ImgUrl, &user.CardNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return nil, err
	}
	return &user, nil
}

func (repo *RepoLayer) DeleteByEmail(ctx context.Context, email string) error {
	row, err := repo.database.ExecContext(ctx, `DELETE FROM "user" WHERE email = $1`, email)
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUserRelation
	}
	return nil
}

func (repo *RepoLayer) Create(ctx context.Context, u *entity.User) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row, err := repo.database.ExecContext(ctx,
		`INSERT INTO "user" (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		u.Name, u.Email, u.Password, timeNow, timeNow)
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUserRelation
	}
	return nil
}

func (repo *RepoLayer) Update(ctx context.Context, uDataChange *entity.User, email string) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row, err := repo.database.ExecContext(ctx,
		`UPDATE "user" SET name = $1, email = $3, phone = $2, img_url = $4, password = $5, card_number = $6, address = $7, updated_at = $8 WHERE email = $9`,
		uDataChange.Name, uDataChange.Email, functions.MaybeNullString(uDataChange.Phone), uDataChange.ImgUrl, uDataChange.Password, functions.MaybeNullString(uDataChange.CardNumber), functions.MaybeNullString(uDataChange.Address), timeNow, email)
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUserRelation
	}
	return nil
}
