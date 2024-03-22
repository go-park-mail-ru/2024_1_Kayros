package user

import (
	"database/sql"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/functions"
)

type UserRepositoryInterface interface {
	GetById(alias.UserId) (*entity.User, error)
	DeleteById(alias.UserId) (bool, error)
	Create(*entity.User) (*entity.User, error)
	Update(*entity.User) (*entity.User, error)
}

type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRepository{
		database: db,
	}
}

func (t *UserRepository) GetById(id alias.UserId) (*entity.User, error) {
	user := &entity.User{}
	row := t.database.QueryRow("SELECT id, name, phone, email, img_url FROM User WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password, &user.ImgUrl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (t *UserRepository) DeleteById(id alias.UserId) (bool, error) {
	row := t.database.QueryRow(`DELETE FROM "User" WHERE id = $1`, id)

	err := row.Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *UserRepository) Create(u *entity.User) (*entity.User, error) {
	hashPassword, err := functions.HashData(u.Password)
	if err != nil {
		return nil, err
	}
	row := t.database.QueryRow(`INSERT INTO "User" (name, phone, email, password, img_url) VALUES ($1, $2, $3, $4, $5)`,
		u.Name, u.Phone, u.Email, hashPassword, u.ImgUrl)

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// пока что полагаю, что валидация будет поддерживать возможные пустые поля
func (t *UserRepository) Update(u *entity.User) (*entity.User, error) {
	user, err := t.GetById(alias.UserId(u.Id))
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

	row := t.database.QueryRow(`UPDATE "User" SET name = $1, phone = $2, email = $3, img_url = $4, password = $5 WHERE id = $6`,
		u.Name, u.Phone, u.Email, u.ImgUrl, u.Password, u.Id)

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (t *UserRepository) CheckPassword(id alias.UserId, password string) (bool, error) {
	hashPassword, err := functions.HashData(password)
	if err != nil {
		return false, err
	}

	user, err := t.GetById(id)
	if err != nil {
		return false, err
	}
	return user.Password == hashPassword, nil
}
