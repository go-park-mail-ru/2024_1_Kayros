package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"2024_1_kayros/internal/delivery/restaurants"
)

const PORT = ":8000"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/restaurants", delivery.RestaurantList).Methods("GET")
	http.Handle("/", r)

	fmt.Printf("server is starting on port 8000\n")
	log.Fatal(http.ListenAndServe(PORT, r))
}
