package delivery

import (
	"database/sql"

	"2024_1_kayros/internal/delivery/middleware"
	"2024_1_kayros/internal/delivery/restaurants"
	"2024_1_kayros/internal/delivery/signin"
	"2024_1_kayros/internal/delivery/signout"
	"2024_1_kayros/internal/delivery/signup"
	"2024_1_kayros/internal/delivery/user"
	"github.com/gorilla/mux"
)

func Setup(db *sql.DB, mux *mux.Router) {
	mux.HandleFunc("/api/v1/signin", signin.SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/api/v1/signup", signup.SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/api/v1/signout", signout.SignOut).Methods("POST").Name("signout")

	mux.HandleFunc("/api/v1/user", user.UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("/api/v1/restaurants", restaurants.RestaurantList).Methods("GET").Name("restaurants")

	handler := middleware.SessionAuthentication(mux)
	handler = middleware.CorsMiddleware(handler)
}
