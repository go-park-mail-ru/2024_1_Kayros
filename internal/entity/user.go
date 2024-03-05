package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// hashData хэширует данные с помощью хэш-функции sha256
func (u *User) hashData(data string) string {
	hashedPassword := sha256.New()
	hashedPassword.Write([]byte(data))
	return hex.EncodeToString(hashedPassword.Sum(nil))
}

// SetPassword устанавливает пароль пользователя
func (u *User) SetPassword(password string) {
	u.Password = u.hashData(password) // возвращает строку
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (u *User) CheckPassword(password string) bool {
	hashPassword := u.hashData(password)
	return u.Password == hashPassword
}

// IsAuthenticated проверяет
func (u *User) IsAuthenticated(r *http.Request) bool {
	userData := r.Context().Value("user")
	return userData != nil
}
