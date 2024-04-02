package auth

import (
	"net/http"

	"2024_1_kayros/internal/repository/session"
	"2024_1_kayros/internal/repository/user"
)

type Usecase interface {
	SignInUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
	SignUpUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
	SignOutUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
}

type UsecaseLayer struct {
	repoUser    user.Repo
	repoSession session.Repo
}

func NewUsecase(repoUserProps user.Repo, repoSessionProps session.Repo) Usecase {
	return &UsecaseLayer{
		repoUser:    repoUserProps,
		repoSession: repoSessionProps,
	}
}
