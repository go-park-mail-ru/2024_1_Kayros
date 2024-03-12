package restaurants

import (
	"encoding/json"
	"net/http"
	"sync"

	"2024_1_kayros/internal/entity"
)

// RestaurantStore хранилище ресторанов
type RestaurantStore struct {
	Data []entity.Restaurant
	sync.RWMutex
}

func (store *RestaurantStore) RestaurantList(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	store.Lock()
	body, err := json.Marshal(store.Data)
	store.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
