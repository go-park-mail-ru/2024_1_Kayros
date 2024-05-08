package repo

import (
	"context"
	"errors"
	"time"

	"2024_1_kayros/internal/utils/myerrors"
	sessionv1 "2024_1_kayros/microservices/session/proto"

	"github.com/redis/go-redis/v9"
)

type Repo interface {
	GetValue(ctx context.Context, key *sessionv1.SessionKey) (*sessionv1.SessionValue, error)
	SetValue(ctx context.Context, data *sessionv1.SetSessionPair) error
	DeleteValue(ctx context.Context, key *sessionv1.SessionKey) error
}

type Layer struct {
	redis *redis.Client
}

func NewLayer(client *redis.Client) Repo {
	return &Layer{
		redis: client,
	}
}

func (repo *Layer) GetValue(ctx context.Context, key *sessionv1.SessionKey) (*sessionv1.SessionValue, error) {
	value, err := repo.redis.Get(ctx, key.GetData()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &sessionv1.SessionValue{}, myerrors.RedisNoData
		}
		return &sessionv1.SessionValue{}, err
	}
	return &sessionv1.SessionValue{Data: value}, nil
}

func (repo *Layer) SetValue(ctx context.Context, data *sessionv1.SetSessionPair) error {
	err := repo.redis.Set(ctx, data.GetKey(), data.GetValue(), 14*24*time.Hour).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

func (repo *Layer) DeleteValue(ctx context.Context, key *sessionv1.SessionKey) error {
	err := repo.redis.Del(ctx, key.GetData()).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}
