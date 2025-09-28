package mq

import (
	config "consumer/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	config     *config.RabbitMQConfig
}

type Message struct {
	UserAgent string
	Token     string
	Message   string
	Receiver  string
}
