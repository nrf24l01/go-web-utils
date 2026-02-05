package rabbitMQ

import (
	"fmt"

	"github.com/nrf24l01/go-web-utils/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RegisterRabbitMQ(cfg *config.RabbitMQConfig) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	rabbitMQ := &RabbitMQ{
		cfg:     cfg,
		Conn:    conn,
		Channel: channel,
	}

	return rabbitMQ, nil
}
