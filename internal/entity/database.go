package entity

import (
	"sync"

	"2024_1_kayros/internal/delivery/restaurants"
	"github.com/satori/uuid"
)

type SystemDatabase struct {
	Users       *UserStore
	Sessions    *SessionStore
	Restaurants *restaurants.RestaurantStore
}

// InitDatabase метод, который инициализирует нашу базу данных
func InitDatabase() *SystemDatabase {
	r := []Restaurant{
		{1, "Пицца 22 см", "Пиццерия с настоящей неаполитанской пиццей", "assets/mocks/restaurants/1.jpg"},
		{2, "Bro&N", "Ресторан классической итальянской кухни", "assets/mocks/restaurants/2.jpg"},
		{3, "#FARШ", "Сеть бургерных с сочным мясом от \"Мираторга\"", "assets/mocks/restaurants/3.jpg"},
		{4, "Loona", "Итальянскую классику в современном прочтении", "assets/mocks/restaurants/4.jpg"},
		{5, "Pino", "Обширное интернациональное меню", "assets/mocks/restaurants/5.jpg"},
		{6, "Sage", "Авторская евпропейская кухня с акцентом на мясо и рыбу", "assets/mocks/restaurants/6.jpg"},
		{7, "TECHNIKUM", "Современное гастробистро с нескучной едой", "assets/mocks/restaurants/7.jpg"},
	}
	rests := restaurants.RestaurantStore{
		Restaurants:      r,
		RestaurantsMutex: sync.RWMutex{},
	}
	users := UserStore{
		Users:      map[DataType]User{},
		UsersMutex: sync.RWMutex{},
	}
	sessions := SessionStore{
		SessionTable:      map[uuid.UUID]DataType{},
		SessionTableMutex: sync.RWMutex{},
	}
	return &SystemDatabase{
		Users:       &users,
		Sessions:    &sessions,
		Restaurants: &rests,
	}
}
