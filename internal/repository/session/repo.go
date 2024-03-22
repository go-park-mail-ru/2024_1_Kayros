package session

import (
	"context"

	"2024_1_kayros/internal/utils/alias"
	"github.com/redis/go-redis/v9"
)

type SessionRepositoryInterface interface {
	GetValue(key alias.SessionKey) (alias.SessionValue, error)
	SetValue(key alias.SessionKey, value alias.SessionValue) error
	DeleteKey(key alias.SessionKey) error
}

type SessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(client *redis.Client) SessionRepositoryInterface {
	return &SessionRepository{
		redis: client,
	}
}

func (t *SessionRepository) GetValue(key alias.SessionKey) (alias.SessionValue, error) {
	value, err := t.redis.Get(context.Background(), string(key)).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	returnValue := alias.SessionValue(value)
	return returnValue, nil
}

func (t *SessionRepository) SetValue(key alias.SessionKey, value alias.SessionValue) error {
	err := t.redis.Set(context.Background(), string(key), string(value), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (t *SessionRepository) DeleteKey(key alias.SessionKey) error {
	err := t.redis.Del(context.Background(), string(key)).Err()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}
