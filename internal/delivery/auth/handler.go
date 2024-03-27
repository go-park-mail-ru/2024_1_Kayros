package auth

import (
	"net/http"

	"2024_1_kayros/internal/usecase/auth"
)

type AuthDelivery struct {
	authUsecase auth.AuthUsecaseInterface
}

func NewAuthDelivery(authUsecase auth.AuthUsecaseInterface) *AuthDelivery {
	return &AuthDelivery{
		authUsecase: authUsecase,
	}
}

func (d *AuthDelivery) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w = d.authUsecase.SignInUser(w, r)
}

func (d *AuthDelivery) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w = d.authUsecase.SignUpUser(w, r)
}

func (d *AuthDelivery) SignOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w = d.authUsecase.SignOutUser(w, r)
}
