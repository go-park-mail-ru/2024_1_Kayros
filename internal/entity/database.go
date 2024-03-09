package entity

import (
	"sync"

	"github.com/satori/uuid"
)

type SystemDatabase struct {
	Users       UserStore    `json:"users"`
	Sessions    SessionStore `json:"sessions"`
	Restaurants []Restaurant `json:"restaurants"`
}

// Allocate метод, который инициализирует нашу базу данных
func (db *SystemDatabase) Allocate() *SystemDatabase {
	restaurants := []Restaurant{
		{1, "Пицца 22 см", "Пиццерия с настоящей неаполитанской пиццей", "assets/mocks/restaurants/1.jpg"},
		{2, "Bro&N", "Ресторан классической итальянской кухни", "assets/mocks/restaurants/2.jpg"},
		{3, "#FARШ", "Сеть бургерных с сочным мясом от \"Мираторга\"", "assets/mocks/restaurants/3.jpg"},
		{4, "Loona", "Итальянскую классику в современном прочтении", "assets/mocks/restaurants/4.jpg"},
		{5, "Pino", "Обширное интернациональное меню", "assets/mocks/restaurants/5.jpg"},
		{6, "Sage", "Авторская евпропейская кухня с акцентом на мясо и рыбу", "assets/mocks/restaurants/6.jpg"},
		{7, "TECHNIKUM", "Современное гастробистро с нескучной едой", "assets/mocks/restaurants/7.jpg"},
	}
	db.Restaurants = restaurants
	db.Users = UserStore{
		Users:      make([]User, 0, 10),
		UsersMutex: sync.RWMutex{},
	}
	db.Sessions = SessionStore{
		SessionTable:      map[uuid.UUID]string{},
		SessionTableMutex: sync.RWMutex{},
	}
	return db
}
