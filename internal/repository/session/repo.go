package session

import (
	"context"
	"time"

	"2024_1_kayros/internal/constants/alias"
	"github.com/redis/go-redis/v9"
)

type SessionRepository interface {
	GetValue(key alias.SessionKey) (alias.SessionValue, error)
	SetValue(key alias.SessionKey, value alias.SessionValue) error
	DeleteKey(key alias.SessionKey) error
}

type SessionTable struct {
	redis *redis.Client
}

func (t *SessionTable) GetValue(key alias.SessionKey) (alias.SessionValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := t.redis.Get(ctx, string(key)).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	returnValue := alias.SessionValue(value)
	return returnValue, nil
}

func (t *SessionTable) SetValue(key alias.SessionKey, value alias.SessionValue) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := t.redis.Set(ctx, string(key), string(value), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (t *SessionTable) DeleteKey(key alias.SessionKey) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := t.redis.Del(ctx, string(key)).Err()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
