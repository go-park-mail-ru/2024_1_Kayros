package entity

import (
	"github.com/satori/uuid"
)

type AuthDatabase struct {
	Users    UserStore
	Sessions SessionStore
}

// InitDatabase метод, который инициализирует нашу базу данных
func InitDatabase() AuthDatabase {
	return AuthDatabase{
		Users: UserStore{
			Data: make(map[string]User),
		},
		Sessions: SessionStore{
			Data: make(map[uuid.UUID]string),
		},
	}
}
