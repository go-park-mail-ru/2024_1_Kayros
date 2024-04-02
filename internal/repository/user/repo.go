package user

import (
	"context"
	"database/sql"
	"errors"

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
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, password, phone, email, address, img_url FROM "User" WHERE id = $1`, uint64(userId))
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password, &user.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *RepoLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, password, phone, email, address, img_url FROM "User" WHERE email = $1`, email)
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password, &user.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *RepoLayer) DeleteById(ctx context.Context, userId alias.UserId) error {
	res, err := repo.database.ExecContext(ctx, `DELETE FROM "User" WHERE id = $1`, uint64(userId))
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return errors.New("Пользователь не был удален")
	}
	return nil
}

func (repo *RepoLayer) DeleteByEmail(ctx context.Context, email string) error {
	res, err := repo.database.ExecContext(ctx, `DELETE FROM "User" WHERE email = $1`, email)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return errors.New("Пользователь не был удален")
	}
	return nil
}

func (repo *RepoLayer) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	hashPassword, err := functions.HashData(u.Password)
	if err != nil {
		return nil, err
	}
	res, err := repo.database.ExecContext(ctx,
		`INSERT INTO "User" (name, phone, email, password, img_url) VALUES ($1, $2, $3, $4, $5)`,
		u.Name, u.Phone, u.Email, hashPassword, u.ImgUrl)

	if err != nil {
		return nil, err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if countRows == 0 {
		return nil, errors.New("Пользователь не был добавлен")
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

	var hashPassword string
	if u.Password == "" {
		hashPassword = user.Password
	} else {
		hashPassword, err = functions.HashData(u.Password)
		if err != nil {
			return nil, err
		}
	}

	res, err := repo.database.ExecContext(ctx,
		`UPDATE "User" SET name = $1, phone = $2, email = $3, img_url = $4, password = $5 WHERE id = $6`,
		u.Name, u.Phone, u.Email, u.ImgUrl, hashPassword, u.Id)

	if err != nil {
		return nil, err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if countRows == 0 {
		return nil, errors.New("Данные о пользователе не были обновлены")
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
