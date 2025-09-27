package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QUEUE = "jobs"
	USERNAME = "guest"
	PASSWORD = "guest"
	HOST     = "localhost"
	PORT     = "5672"
	SENDMESSAGE = "message-send"
	SERVERPORT = ":8000"
)

type MQ struct {
	connection *amqp.Connection
	channel *amqp.Channel
}

type Message struct {
	UserAgent string
	Token string
	Message string
	Reciever string
}