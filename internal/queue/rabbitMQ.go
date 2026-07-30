package mq

import (
	config "consumer/internal/config"

	"encoding/json"
	"errors"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Create(cfg *config.RabbitMQConfig) (*MQ, error) {
	connection, err := amqp.Dial(cfg.GetRabbitMQURL())
	if err != nil {
		return nil, err
	}
	
	channel, err :=  connection.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		cfg.DeadQueue,
		true,	// durable
		false,	// autoDeleted
		false,	// exclusive
		false,	// noWait
		amqp.Table{
			"x-dead-letter-exchange": "",
			"x-dead-letter-routing-key": cfg.Queue,
			"x-message-ttl": int32(1000),
    	},
	)
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		cfg.Queue,
		true,	// durable
		false,	// autoDeleted
		false,	// exclusive
		false,	// noWait
		amqp.Table{
			"x-dead-letter-exchange": "",
			"x-dead-letter-routing-key": cfg.DeadQueue,
		},
	)
	if err != nil {
		return nil, err
	}

	return &MQ{
		connection, 
		channel,
		cfg,
	}, nil
}

func (mq *MQ) Publish(message Message) error {
	if mq.connection == nil || mq.channel == nil{
		return errors.New("no channel found") 
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return mq.channel.Publish(
		"",
		mq.config.Queue,
		false,
		false,
		amqp.Publishing{
			MessageId: uuid.NewString(),
			ContentType: "application/json",
			Body: body,
			Type: mq.config.Type,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (mq *MQ) Consume() (<-chan amqp.Delivery, error) {
	if mq.connection == nil || mq.channel == nil{
		return nil, errors.New("no channel found") 
	}
	
	return mq.channel.Consume(
		mq.config.Queue,
		"",
		false,	// autoAck
		false,	// exclusive
		false,	// noLocal
		false,	// noWait
		nil,	// args
	)
}

func (mq *MQ) Close() error{
	if  channelErr := mq.channel.Close(); 
		channelErr != nil {
		return channelErr
	}
	
	if  connectionErr := mq.connection.Close(); 
		connectionErr != nil {
		return connectionErr
	}

	return nil
}