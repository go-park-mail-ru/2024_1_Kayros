package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"sync"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HashData хэширует данные с помощью хэш-функции sha256
func HashData(data string) string {
	hashedPassword := sha256.New()
	hashedPassword.Write([]byte(data))
	return hex.EncodeToString(hashedPassword.Sum(nil))
}

// SetPassword устанавливает пароль пользователя
func (u *User) SetPassword(password string) (DataType, ErrorType) {
	regexPassword := regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	if regexPassword.MatchString(password) {
		u.Password = HashData(password) // возвращает строку
		return GenerateResponse(true)
	}
	return RaiseError("Предоставлены неверные учетные данные")
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (u *User) CheckPassword(password string) bool {
	hashPassword := HashData(password)
	return u.Password == hashPassword
}

// UserStore хранилище с пользователями
type UserStore struct {
	Users      map[string]User
	UsersMutex sync.RWMutex
}

// GetUser возвращает пользователя
func (s *UserStore) GetUser(field string) (DataType, ErrorType) {
	s.UsersMutex.RLock()
	user, userExist := s.Users[field]
	s.UsersMutex.RUnlock()
	if userExist {
		return GenerateResponse(user)
	}
	return RaiseError("Предоставлены неверные учетные данные")
}

// SetNewUser добавляет нового пользователя в БД
func (s *UserStore) SetNewUser(field string, data User) (DataType, ErrorType) {
	// пока что мы проверяем по почте
	regexEmail := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if regexEmail.MatchString(field) {
		s.UsersMutex.Lock()
		s.Users[field] = data
		s.Users[field].SetPassword(data.Password)
		s.UsersMutex.Unlock()
		return GenerateResponse(data)
	}
	return RaiseError("Предоставлены неверные учетные данные")
}
