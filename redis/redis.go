package redis

import (
	"context"
	"log"

	"github.com/nrf24l01/go-web-utils/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(config *config.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	return &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}
}
