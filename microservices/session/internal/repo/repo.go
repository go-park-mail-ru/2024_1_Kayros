package repo

import (
	"context"
	"errors"
	"time"

	"2024_1_kayros/gen/go/session"
	metrics "2024_1_kayros/microservices/metrics"

	"2024_1_kayros/internal/utils/myerrors"

	"github.com/redis/go-redis/v9"
)

type Repo interface {
	GetValue(ctx context.Context, key string) (*session.SessionValue, error)
	SetValue(ctx context.Context, key string, value string) error
	DeleteValue(ctx context.Context, key string) error
}

type Layer struct {
	redis   *redis.Client
	metrics *metrics.MicroserviceMetrics
}

func NewLayer(client *redis.Client, metrics *metrics.MicroserviceMetrics) Repo {
	return &Layer{
		redis:   client,
		metrics: metrics,
	}
}

func (repo *Layer) GetValue(ctx context.Context, key string) (*session.SessionValue, error) {
	timeNow := time.Now()
	value, err := repo.redis.Get(ctx, key).Result()
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.REDIS).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &session.SessionValue{}, myerrors.RedisNoData
		}
		return &session.SessionValue{}, err
	}
	return &session.SessionValue{Data: value}, nil
}

func (repo *Layer) SetValue(ctx context.Context, key string, value string) error {
	timeNow := time.Now()
	err := repo.redis.Set(ctx, key, value, 14*24*time.Hour).Err()
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.REDIS).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

func (repo *Layer) DeleteValue(ctx context.Context, key string) error {
	timeNow := time.Now()
	err := repo.redis.Del(ctx, key).Err()
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.REDIS).Observe(float64(timeEnd.Milliseconds()))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}
