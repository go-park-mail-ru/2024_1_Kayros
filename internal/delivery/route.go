package delivery

import (
	"database/sql"

	"2024_1_kayros/internal/delivery/middleware"
	"github.com/gorilla/mux"
)

func Setup(db *sql.DB, mux *mux.Router) {
	mux.HandleFunc("/api/v1/signin", SignIn).Methods("POST").Name("signin")
	mux.HandleFunc("/api/v1/signup", SignUp).Methods("POST").Name("signup")
	mux.HandleFunc("/api/v1/signout", SignOut).Methods("POST").Name("signout")

	mux.HandleFunc("/api/v1/user", UserData).Methods("GET").Name("userdata")
	mux.HandleFunc("/api/v1/restaurants", RestaurantList).Methods("GET").Name("restaurants")

	handler := middleware.SessionAuthentication(mux)
	handler = middleware.CorsMiddleware(handler)
}
