package session

import (
	"context"
	"errors"
	"time"

	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
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
	methodName := cnst.NameMethodGetValue
	value, err := repo.redis.Get(ctx, string(key)).Result()
	if err == redis.Nil {
		err = errors.New("Такого ключа в Redis не существует")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return "", nil
	} else if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return "", err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return alias.SessionValue(value), nil
}

func (repo *RepoLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, requestId string) error {
	methodName := cnst.NameMethodSetValue
	err := repo.redis.Set(ctx, string(key), string(value), 14*24*time.Hour).Err()
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) DeleteKey(ctx context.Context, key alias.SessionKey, requestId string) (bool, error) {
	methodName := cnst.NameMethodDeleteKey
	err := repo.redis.Del(ctx, string(key)).Err()
	if errors.Is(err, redis.Nil) {
		err = errors.New("Такого ключа в Redis не существует")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return false, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return false, err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return true, nil
}
