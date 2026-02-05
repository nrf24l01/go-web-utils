package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT" envDefault:"http://127.0.0.1:9000"`
	AccessKey string `env:"S3_ACCESS_KEY" envDefault:""`
	SecretKey string `env:"S3_SECRET_KEY" envDefault:""`
	UseSSL    bool   `env:"S3_USE_SSL" envDefault:"true"`
	BaseURL   string `env:"S3_BASE_URL" envDefault:"http://127.0.0.1:9000"`
}

func LoadS3ConfigFromEnv() *S3Config {
	config := &S3Config{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}
