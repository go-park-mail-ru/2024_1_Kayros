package entity

import (
	"errors"
	"sync"

	"github.com/satori/uuid"
)

// SessionStore хранилище сессий пользователей
type SessionStore struct {
	SessionTable map[uuid.UUID]string // ключ - сессия, значение - либо почта, либо номер телефона
	sync.RWMutex
}

func (s *SessionStore) GetValue(key uuid.UUID) (string, error) {
	s.RLock()
	email, emailExist := s.SessionTable[key]
	s.RUnlock()

	if emailExist {
		return email, nil
	}
	return "", errors.New(BadPermission)
}

func (s *SessionStore) SetNewSession(email string) {
	sessionId := uuid.NewV4()
	s.Lock()
	s.SessionTable[sessionId] = email
	s.Unlock()
}

func (s *SessionStore) HasKey(key uuid.UUID) bool {
	s.RLock()
	_, hasKey := s.SessionTable[key]
	s.RUnlock()
	return hasKey
}
