package restaurants

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"2024_1_kayros/internal/entity"
)

// RestaurantStore хранилище ресторанов
type RestaurantStore struct {
	Restaurants []entity.Restaurant
	sync.RWMutex
}

func (store *RestaurantStore) RestaurantList(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	store.Lock()
	body, err := json.Marshal(store.Restaurants)
	store.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		entity.RaiseError("Ошибка формирования ответа")
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
