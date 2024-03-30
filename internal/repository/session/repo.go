package session

import (
	"context"

	"2024_1_kayros/internal/utils/alias"
	"github.com/redis/go-redis/v9"
)

type Repo interface {
	GetValue(key alias.SessionKey) (alias.SessionValue, error)
	SetValue(key alias.SessionKey, value alias.SessionValue) (bool, error)
	DeleteKey(key alias.SessionKey) (bool, error)
}

type RepoLayer struct {
	redis *redis.Client
}

func NewRepoLayer(client *redis.Client) Repo {
	return &RepoLayer{
		redis: client,
	}
}

func (t *RepoLayer) GetValue(key alias.SessionKey) (alias.SessionValue, error) {
	value, err := t.redis.Get(context.Background(), string(key)).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	returnValue := alias.SessionValue(value)
	return returnValue, nil
}

func (t *RepoLayer) SetValue(key alias.SessionKey, value alias.SessionValue) (bool, error) {
	err := t.redis.Set(context.Background(), string(key), string(value), 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *RepoLayer) DeleteKey(key alias.SessionKey) (bool, error) {
	err := t.redis.Del(context.Background(), string(key)).Err()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
