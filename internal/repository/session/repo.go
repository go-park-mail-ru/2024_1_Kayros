package session

import (
	"context"

	"2024_1_kayros/internal/utils/alias"
	"github.com/redis/go-redis/v9"
)

type Repo interface {
	GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) (bool, error)
	DeleteKey(ctx context.Context, key alias.SessionKey) (bool, error)
}

type RepoLayer struct {
	redis *redis.Client
}

func NewRepoLayer(client *redis.Client) Repo {
	return &RepoLayer{
		redis: client,
	}
}

func (repo *RepoLayer) GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error) {
	value, err := repo.redis.Get(ctx, string(key)).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	returnValue := alias.SessionValue(value)
	return returnValue, nil
}

func (repo *RepoLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) (bool, error) {
	err := repo.redis.Set(ctx, string(key), string(value), 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *RepoLayer) DeleteKey(ctx context.Context, key alias.SessionKey) (bool, error) {
	err := repo.redis.Del(ctx, string(key)).Err()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
