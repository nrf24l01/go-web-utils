package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type WebAppConfig struct {
	AppHost     string `env:"APP_HOST" envDefault:"8080"`
	AllowOrigin string `env:"ALLOW_ORIGIN" envDefault:"http://127.0.0.1:5137"`
}

func LoadWebAppConfigFromEnv() *WebAppConfig {
	config := &WebAppConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}