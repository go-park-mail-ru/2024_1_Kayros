package redis

import (
	"context"
	"fmt"
	"time"

	"2024_1_kayros/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Init(cfg *config.Project, logger *zap.Logger, dbNum int) *redis.Client {
	cfgRedis := cfg.Redis
	redisAddress := fmt.Sprintf("%s:%d", cfgRedis.Host, cfgRedis.Port)
	r := redis.NewClient(&redis.Options{
		DB:   dbNum,
		Addr: redisAddress,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Не удалось подключиться к Redis", zap.Error(err))
	}

	logger.Info("Redis успешно подключен")
	return r
}
