package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type TgWebAppConfig struct {
	TgBotToken          string `env:"TG_BOT_TOKEN" envDefault:""`
	InitDataExpireHours int    `env:"INIT_DATA_EXPIRE_HOURS" envDefault:"24"`
}

func LoadTgWebAppConfigFromEnv() *TgWebAppConfig {
	config := &TgWebAppConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}
