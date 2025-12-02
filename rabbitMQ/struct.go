package rabbitMQ

import (
	"github.com/nrf24l01/go-web-utils/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	cfg *config.RabbitMQConfig
	Conn *amqp.Connection
	Channel *amqp.Channel
}