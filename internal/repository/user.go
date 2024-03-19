package repository

import "database/sql"

type UserRepoInterface interface {
}

type UserRepo struct {
	database *sql.DB
}

func GetUserRepo(db *sql.DB) UserRepoInterface {
	return &UserRepo{
		database: db,
	}
}

func (database *UserRepo) GetByEmail(email string) {

}
