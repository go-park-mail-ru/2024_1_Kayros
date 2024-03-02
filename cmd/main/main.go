package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/restaurants", RestaurantList)
	r.HandleFunc("/sigin", SigIn)
	r.HandleFunc("/signup", SignUp)
	r.HandleFunc("/signout", SignOut)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
