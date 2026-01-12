package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type PGConfig struct {
	PGHost     string `env:"POSTGRES_HOST"`
	PGPort     string `env:"POSTGRES_PORT"`
	PGUser     string `env:"POSTGRES_USER"`
	PGPassword string `env:"POSTGRES_PASSWORD"`
	PGDatabase string `env:"POSTGRES_DB"`
	PGSSLMode  string `env:"POSTGRES_SSLMODE"`
	PGTimeZone string `env:"POSTGRES_TIMEZONE"`
}

func LoadPGConfigFromEnv() *PGConfig {
	config := &PGConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}