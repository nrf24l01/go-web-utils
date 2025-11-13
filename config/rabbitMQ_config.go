package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type RabbitMQConfig struct {
	RabbitMQHost     string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	RabbitMQPort     string `env:"RABBITMQ_PORT" envDefault:"5672"`
	RabbitMQUser     string `env:"RABBITMQ_USER" envDefault:"guest"`
	RabbitMQPassword string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
	RabbitMQVHost    string `env:"RABBITMQ_VHOST" envDefault:"/"`
}

func LoadRabbitMQConfigFromEnv() *RabbitMQConfig {
	config := &RabbitMQConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return config
}