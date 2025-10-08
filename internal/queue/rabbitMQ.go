package mq

import (
	config "consumer/internal/config"
	
	"encoding/json"
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

func Create(cfg *config.RabbitMQConfig) (*MQ, error) {
	connection, err := amqp091.Dial(cfg.GetRabbitMQURL())
	if err != nil {
		return nil, err
	}
	
	channel, err :=  connection.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
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

		amqp091.Publishing{
			ContentType: "application/json",
			Body: body,
			Type: mq.config.Type,
		},
	)
}

func (mq *MQ) Consume() (<-chan amqp091.Delivery, error) {
	if mq.connection == nil || mq.channel == nil{
		return nil, errors.New("no channel found") 
	}
	
	return mq.channel.Consume(
		mq.config.Queue,
		"",	
		false,
		false,
		false,
		false,
		nil,
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