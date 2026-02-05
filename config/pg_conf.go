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
	Migrations string `env:"POSTGRES_MIGRATIONS_DIR"`
}

func LoadPGConfigFromEnv() *PGConfig {
	config := &PGConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}

func (cfg *PGConfig) GetDSN() string {
	dsn := "postgres://"
	dsn += cfg.PGUser
	if cfg.PGPassword != "" {
		dsn += ":" + cfg.PGPassword
	}
	dsn += "@" + cfg.PGHost
	if cfg.PGPort != "" {
		dsn += ":" + cfg.PGPort
	}
	dsn += "/" + cfg.PGDatabase
	params := ""
	if cfg.PGSSLMode != "" {
		params += "sslmode=" + cfg.PGSSLMode
	}
	if cfg.PGTimeZone != "" {
		if params != "" {
			params += "&"
		}
		params += "timezone=" + cfg.PGTimeZone
	}
	if params != "" {
		dsn += "?" + params
	}
	return dsn
}
