package entity

import (
	"github.com/satori/uuid"
)

type Database struct {
	SessionTable map[uuid.UUID]string // ключ - сессия, значение - идентификатор пользователя
	Users        map[string]User      // ключ - почта пользователя, значение - данные пользователя (экземпляр структуры)
}

func (db *Database) InitializeDatabase() {
	db.setInitialData()
}

func (db *Database) setInitialData() {
	users := []User{
		{Id: 1, Name: "Ivan", Email: "ivan@yandex.ru"},
		{Id: 2, Name: "Sofia", Email: "sofia@yandex.ru"},
		{Id: 3, Name: "Bogdan", Email: "bogdan@yandex.ru"},
		{Id: 4, Name: "Pasha", Email: "pasha@yandex.ru"},
		{Id: 5, Name: "Ilya", Email: "ilya@yandex.ru"},
	}

	for _, user := range users {
		db.Users[user.Email] = user
	}
}
