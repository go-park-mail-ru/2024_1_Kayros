package user

import (
	"context"
	"database/sql"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/functions"
)

// Передаем контекст запроса пользователя (! возможно лучше еще переопределить контекстом WithTimeout)
type Repo interface {
	GetById(ctx context.Context, userId alias.UserId) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	DeleteById(ctx context.Context, userId alias.UserId) error
	DeleteByEmail(ctx context.Context, email string) error

	Create(ctx context.Context, u *entity.User) (*entity.User, error)
	Update(ctx context.Context, u *entity.User) (*entity.User, error)

	IsExistById(ctx context.Context, userId alias.UserId) (bool, error)
	IsExistByEmail(ctx context.Context, email string) (bool, error)

	CheckPassword(ctx context.Context, email string, password string) (bool, error)
}

type RepoLayer struct {
	database *sql.DB
}

func NewRepoLayer(db *sql.DB) Repo {
	return &RepoLayer{
		database: db,
	}
}

func (repo *RepoLayer) GetById(ctx context.Context, userId alias.UserId) (*entity.User, error) {
	user := &entity.User{}
	row := repo.database.QueryRowContext(ctx, `SELECT id, name, phone, email, img_url FROM "User" WHERE id = $1`, uint64(userId))
	err := row.Scan(user.Id, user.Name, user.Phone, user.Email, user.Password, user.ImgUrl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *RepoLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	row := repo.database.QueryRowContext(ctx, `SELECT id, name, phone, email, password, img_url FROM "User" WHERE email = $1`, email)

	err := row.Scan(user.Id, user.Name, user.Phone, user.Email, user.Password, user.ImgUrl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *RepoLayer) DeleteById(ctx context.Context, userId alias.UserId) error {
	row := repo.database.QueryRowContext(ctx, `DELETE FROM "User" WHERE id = $1`, uint64(userId))

	err := row.Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RepoLayer) DeleteByEmail(ctx context.Context, email string) error {
	row := repo.database.QueryRowContext(ctx, `DELETE FROM "User" WHERE email = $1`, email)

	err := row.Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RepoLayer) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	hashPassword, err := functions.HashData(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashPassword
	row := repo.database.QueryRowContext(ctx, `INSERT INTO "User" (name, phone, email, password, img_url) VALUES ($1, $2, $3, $4, $5)`,
		u.Name, u.Phone, u.Email, u.Password, u.ImgUrl)

	err = row.Err()
	if err != nil {
		return nil, err
	}

	user, err := repo.GetByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *RepoLayer) Update(ctx context.Context, u *entity.User) (*entity.User, error) {
	user, err := repo.GetById(ctx, alias.UserId(u.Id))
	if err != nil {
		return nil, err
	}

	if u.Password == "" {
		u.Password = user.Password
	} else {
		u.Password, err = functions.HashData(u.Password)
		if err != nil {
			return nil, err
		}
	}

	row := repo.database.QueryRowContext(ctx, `UPDATE "User" SET name = $1, phone = $2, email = $3, img_url = $4, password = $5 WHERE id = $6`,
		u.Name, u.Phone, u.Email, u.ImgUrl, u.Password, u.Id)

	err = row.Err()
	if err != nil {
		return nil, err
	}

	user, err = repo.GetById(ctx, alias.UserId(u.Id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (repo *RepoLayer) CheckPassword(ctx context.Context, email string, password string) (bool, error) {
	hashPassword, err := functions.HashData(password)
	if err != nil {
		return false, err
	}

	user, err := repo.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user.Password == hashPassword, nil
}

func (repo *RepoLayer) IsExistById(ctx context.Context, userId alias.UserId) (bool, error) {
	_, err := repo.GetById(ctx, userId)
	// нужно узнать, выдаст ли оштбку отсутсвтие записи в БД
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *RepoLayer) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	_, err := repo.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return true, nil
}
