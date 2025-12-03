package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type JWTConfig struct {
	AccessJWTSecret  string `env:"ACCESS_JWT_SECRET"`
	RefreshJWTSecret string `env:"REFRESH_JWT_SECRET"`

	AccessTokenExpiryMinutes  int `env:"ACCESS_TOKEN_EXPIRY_MINUTES" envDefault:"15"`
	RefreshTokenExpiryMinutes int `env:"REFRESH_TOKEN_EXPIRY_MINUTES" envDefault:"10080"` // 7 days
}

func LoadJWTConfigFromEnv() *JWTConfig {
	config := &JWTConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}