package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type PGConfig struct {
	PGHost     string `env:"PG_HOST" envDefault:"localhost"`
	PGPort     string `env:"PG_PORT" envDefault:"5432"`
	PGUser     string `env:"PG_USER" envDefault:"postgres"`
	PGPassword string `env:"PG_PASSWORD" envDefault:"password"`
	PGDatabase string `env:"PG_DATABASE" envDefault:"postgres"`
	PGSSLMode  string `env:"PG_SSLMODE" envDefault:"disable"`
	PGTimeZone string `env:"PG_TIMEZONE" envDefault:"UTC"`
}

func LoadPGConfigFromEnv() *PGConfig {
	config := &PGConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}