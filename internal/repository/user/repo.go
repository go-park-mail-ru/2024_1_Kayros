package user

import "database/sql"

type UserRepoInterface interface {
	func GetByEmail (email string)
}

type UserRepo struct {
	database *sql.DB
}

func GetUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		database: db,
	}
}

func (repo *UserRepo) GetByEmail(email string) {
	repo
}
