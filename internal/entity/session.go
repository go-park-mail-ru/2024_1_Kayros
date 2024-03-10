package entity

import (
	"regexp"
	"sync"

	"github.com/satori/uuid"
)

// SessionStore хранилище сессий пользователей
type SessionStore struct {
	SessionTable      map[uuid.UUID]DataType // ключ - сессия, значение - идентификатор пользователя
	SessionTableMutex sync.RWMutex
}

func (s *SessionStore) GetValue(key uuid.UUID) (DataType, ErrorType) {
	s.SessionTableMutex.RLock()
	email, emailExist := s.SessionTable[key]
	s.SessionTableMutex.RUnlock()
	if emailExist {
		return GenerateResponse(email)
	}
	return RaiseError("Ошибка получения email")
}

func (s *SessionStore) SetNewSession(email string) (DataType, ErrorType) {
	regexEmail := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if regexEmail.MatchString(email) {
		sessionId := uuid.NewV4()
		s.SessionTableMutex.Lock()
		s.SessionTable[sessionId] = email
		s.SessionTableMutex.Unlock()
		return GenerateResponse(sessionId)
	} else {
		return RaiseError("Предоставлены неверные учетные данные")
	}
}

func (s *SessionStore) HasKey(key uuid.UUID) (DataType, ErrorType) {
	s.SessionTableMutex.RLock()
	_, hasKey := s.SessionTable[key]
	s.SessionTableMutex.RUnlock()
	if hasKey {
		return GenerateResponse(true)
	} else {
		return GenerateResponse(false)
	}
}
