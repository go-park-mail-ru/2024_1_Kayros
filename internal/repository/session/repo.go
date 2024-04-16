package session

import (
	"context"
	"errors"
	"time"

	"2024_1_kayros/internal/utils/alias"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Repo interface {
	GetValue(ctx context.Context, key alias.SessionKey, requestId string) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, requestId string) error
	DeleteKey(ctx context.Context, key alias.SessionKey, requestId string) (bool, error)
}

type RepoLayer struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewRepoLayer(client *redis.Client, loggerProps *zap.Logger) Repo {
	return &RepoLayer{
		redis:  client,
		logger: loggerProps,
	}
}

func (repo *RepoLayer) GetValue(ctx context.Context, key alias.SessionKey, requestId string) (alias.SessionValue, error) {
	value, err := repo.redis.Get(ctx, string(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return alias.SessionValue(value), nil
}

func (repo *RepoLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, requestId string) error {
	err := repo.redis.Set(ctx, string(key), string(value), 14*24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (repo *RepoLayer) DeleteKey(ctx context.Context, key alias.SessionKey, requestId string) (bool, error) {
	err := repo.redis.Del(ctx, string(key)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
