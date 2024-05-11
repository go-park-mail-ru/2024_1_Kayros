package repo

import (
	"context"
	"errors"
	"time"

	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/internal/utils/myerrors"

	"github.com/redis/go-redis/v9"
)

type Repo interface {
	GetValue(ctx context.Context, key string) (*session.SessionValue, error)
	SetValue(ctx context.Context, key string, value string) error
	DeleteValue(ctx context.Context, key string) error
}

type Layer struct {
	redis *redis.Client
}

func NewLayer(client *redis.Client) Repo {
	return &Layer{
		redis: client,
	}
}

func (repo *Layer) GetValue(ctx context.Context, key string) (*session.SessionValue, error) {
	value, err := repo.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &session.SessionValue{}, myerrors.RedisNoData
		}
		return &session.SessionValue{}, err
	}
	return &session.SessionValue{Data: value}, nil
}

func (repo *Layer) SetValue(ctx context.Context, key string, value string) error {
	err := repo.redis.Set(ctx, key, value, 14*24*time.Hour).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

func (repo *Layer) DeleteValue(ctx context.Context, key string) error {
	err := repo.redis.Del(ctx, key).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}
