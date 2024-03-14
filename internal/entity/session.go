package entity

import (
	"errors"
	"sync"

	"github.com/satori/uuid"
)

// SessionStore хранилище сессий пользователей
type SessionStore struct {
	Data map[uuid.UUID]string // ключ - сессия, значение - либо почта, либо номер телефона
	sync.RWMutex
}

func (s *SessionStore) GetValueByKey(key uuid.UUID) (string, error) {
	s.RLock()
	email, emailExist := s.Data[key]
	s.RUnlock()

	if emailExist {
		return email, nil
	}
	return "", errors.New(BadPermission)
}

func (s *SessionStore) SetNewSession(sessionId uuid.UUID, email string) {
	s.Lock()
	s.Data[sessionId] = email
	s.Unlock()
}

func (s *SessionStore) DeleteSession(sessionId uuid.UUID) {
	s.Lock()
	delete(s.Data, sessionId)
	s.Unlock()
}

func (s *SessionStore) HasKey(key uuid.UUID) bool {
	s.RLock()
	_, hasKey := s.Data[key]
	s.RUnlock()
	return hasKey
}
