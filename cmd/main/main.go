package main

import (
	"2024_1_kayros/internal/delivery/restaurants"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/restaurants", delivery.RestaurantList).Methods("GET")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", r))
}
