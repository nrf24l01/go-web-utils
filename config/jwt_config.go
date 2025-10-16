package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type JWTConfig struct {
	AccessJWTSecret  string `env:"ACCESS_JWT_SECRET"`
	RefreshJWTSecret string `env:"REFRESH_JWT_SECRET"`
}

func LoadJWTConfigFromEnv() *JWTConfig {
	config := &JWTConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}