package entity

import (
	"github.com/satori/uuid"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSession struct {
	Email     string
	SessionId uuid.UUID
}
