package session

import (
	"context"
	"errors"
	"time"

	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/redis/go-redis/v9"
)

type Repo interface {
	GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) error
	DeleteKey(ctx context.Context, key alias.SessionKey) error
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
	if err != nil {
		if err == redis.Nil {
			return "", myerrors.RedisNoData
		}
		return "", err
	}
	return alias.SessionValue(value), nil
}

func (repo *RepoLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) error {
	err := repo.redis.Set(ctx, string(key), string(value), 14*24*time.Hour).Err()
	if err != nil {
		if err == redis.Nil {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

func (repo *RepoLayer) DeleteKey(ctx context.Context, key alias.SessionKey) error {
	err := repo.redis.Del(ctx, string(key)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}
