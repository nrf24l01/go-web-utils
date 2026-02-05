package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Argon2idConfig struct {
	Memory      uint32 `env:"ARGON2ID_MEMORY" envDefault:"131072"` // in KB
	Time        uint32 `env:"ARGON2ID_TIME" envDefault:"6"`
	Parallelism uint8  `env:"ARGON2ID_PARALLELISM" envDefault:"4"`
	SaltLength  uint32 `env:"ARGON2ID_SALT_LENGTH" envDefault:"16"`
	KeyLength   uint32 `env:"ARGON2ID_KEY_LENGTH" envDefault:"32"`
}

func LoadArgon2idConfigFromEnv() *Argon2idConfig {
	config := &Argon2idConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse Argon2id environment variables: %v", err)
	}
	return config
}
