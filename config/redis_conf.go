package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type RedisConfig struct {
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
}

func LoadRedisConfigFromEnv() *RedisConfig {
	config := &RedisConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}