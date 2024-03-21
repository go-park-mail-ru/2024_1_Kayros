package user

import (
	"database/sql"

	"2024_1_kayros/internal/entity"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetByEmail(email string)
}

type UserTable struct {
	database *sql.DB
	redis    *redis.Client
}

func (t *UserTable) GetUserById() *entity.User {

}

func (repo *UserRepo) GetByEmail(email string) {
	repo
}
