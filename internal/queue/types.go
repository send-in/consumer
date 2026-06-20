package mq

import (
	config "consumer/internal/config"

	ampq "github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	connection *ampq.Connection
	channel    *ampq.Channel
	config     *config.RabbitMQConfig
}

type Message struct {
	UserAgent string
	JSession  string
	Token     string

	Message   string
	Receiver  string
}
