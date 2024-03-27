package redis

import (
	"context"
	"fmt"
	"time"

	"2024_1_kayros/config"
	"github.com/redis/go-redis/v9"
)

func RedisInit(cfg *config.Project) (*redis.Client, error) {
	cfgRedis := cfg.Redis
	redisAddress := fmt.Sprintf("%s:%d", cfgRedis.Host, cfgRedis.Port)
	r := redis.NewClient(&redis.Options{
		DB:       cfgRedis.Database,
		Addr:     redisAddress,
		Username: cfgRedis.User,
		Password: cfgRedis.Password,
	})

	const maxPingTime = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), maxPingTime)
	defer cancel()

	_, err := r.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return r, nil
}
