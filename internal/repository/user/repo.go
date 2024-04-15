package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	cnst "2024_1_kayros/internal/utils/constants"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
)

type Repo interface {
	GetById(ctx context.Context, userId alias.UserId) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	DeleteById(ctx context.Context, userId alias.UserId) error
	DeleteByEmail(ctx context.Context, email string) error

	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, uDataChange *entity.User, email string) error
}

type RepoLayer struct {
	database *sql.DB
	logger   *zap.Logger
}

func NewRepoLayer(db *sql.DB, loggerProps *zap.Logger) Repo {
	return &RepoLayer{
		database: db,
		logger:   loggerProps,
	}
}

func (repo *RepoLayer) GetById(ctx context.Context, userId alias.UserId) (*entity.User, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, email, COALESCE(phone, ''), password, COALESCE(address, ''), img_url, COALESCE(card_number, '')  FROM "user" WHERE id = $1`, uint64(userId))
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Address, &user.ImgUrl, &user.CardNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *RepoLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, email, COALESCE(phone, ''), password, COALESCE(address, ''), img_url, COALESCE(card_number, '')  FROM "user" WHERE email = $1`, email)
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Address, &user.ImgUrl, &user.CardNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *RepoLayer) DeleteById(ctx context.Context, userId alias.UserId) error {
	row := repo.database.QueryRowContext(ctx, `DELETE FROM "user" WHERE id = $1 RETURNING id, email`, uint64(userId))
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func (repo *RepoLayer) DeleteByEmail(ctx context.Context, email string) error {
	row := repo.database.QueryRowContext(ctx, `DELETE FROM "user" WHERE email = $1 RETURNING id, email`, email)
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func (repo *RepoLayer) Create(ctx context.Context, u *entity.User) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row := repo.database.QueryRowContext(ctx,
		`INSERT INTO "user" (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, email`,
		u.Name, u.Email, u.Password, timeNow, timeNow)
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}
	return nil
}

func (repo *RepoLayer) Update(ctx context.Context, uDataChange *entity.User, email string) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row := repo.database.QueryRowContext(ctx,
		`UPDATE "user" SET name = $1, phone = $2, email = $3, img_url = $4, password = $5, card_number = $6, address = $7, updated_at = $8 WHERE email = $9 RETURNING id, email`,
		uDataChange.Name, uDataChange.Phone, uDataChange.Email, uDataChange.ImgUrl, uDataChange.Password, uDataChange.CardNumber, uDataChange.Address, timeNow, email)
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if strings.Contains(err.Error(), "user_email_key") {
			return err
		}
		return err
	}
	return nil
}
