package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"2024_1_kayros/internal/entity"
)

type RestaurantStore struct {
	rests []*entity.Restaurant
	mu    sync.RWMutex
}

func NewRestaurantStore() *RestaurantStore {
	return &RestaurantStore{
		rests: []*entity.Restaurant{
			{1, "Пицца 22 см", "Пиццерия с настоящей неаполитанской пиццей", "assets/mocks/restaurants/1.jpg"},
			{2, "Bro&N", "Ресторан классической итальянской кухни", "assets/mocks/restaurants/2.jpg"},
			{3, "#FARШ", "Сеть бургерных с сочным мясом от \"Мираторга\"", "assets/mocks/restaurants/3.jpg"},
			{4, "Loona", "Итальянскую классику в современном прочтении", "assets/mocks/restaurants/4.jpg"},
			{5, "Pino", "Обширное интернациональное меню", "assets/mocks/restaurants/5.jpg"},
			{6, "Sage", "Авторская евпропейская кухня с акцентом на мясо и рыбу", "assets/mocks/restaurants/6.jpg"},
			{7, "TECHNIKUM", "Современное гастробистро с нескучной едой", "assets/mocks/restaurants/7.jpg"},
		},
		mu: sync.RWMutex{},
	}
}

func (store *RestaurantStore) RestaurantList(w http.ResponseWriter, _ *http.Request) {
	store.mu.Lock()
	body, err := json.Marshal(store.rests)
	store.mu.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	_, err = w.Write(body)
	if err != nil {
		fmt.Printf("Write failed: %v", err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
