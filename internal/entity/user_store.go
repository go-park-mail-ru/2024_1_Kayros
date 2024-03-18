package entity

import (
	"errors"
	"sync"
)

// UserStore хранилище с пользователями
type UserStore struct {
	Data map[string]User
	sync.RWMutex
}

// GetUser возвращает пользователя
func (s *UserStore) GetUser(key string) (User, error) {
	s.RLock()
	user, userExist := s.Data[key]
	s.RUnlock()
	if userExist {
		return user, nil
	}
	return User{}, errors.New(BadAuthCredentials)
}

// SetNewUser добавляет нового пользователя в БД
func (s *UserStore) SetNewUser(key string, data User) (bool, error) {
	s.Lock()
	s.Data[key] = data
	s.Unlock()

	s.RLock()
	_, userExist := s.Data[key]
	s.RUnlock()

	if userExist {
		return true, nil
	}
	return false, errors.New(BadRegCredentials)
}
