package mq

import (
	"consumer/util"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateConnection(url string) (*MQ, error) {

	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	
	channel, err :=  connection.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		QUEUE,
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
		QUEUE,
		false,
		false,

		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
			Type: SENDMESSAGE,
		},
	)
}

func (mq *MQ) Consume() (<-chan amqp.Delivery, error) {

	if mq.connection == nil || mq.channel == nil{
		return nil, errors.New("no channel found") 
	}
	
	return mq.channel.Consume(
		QUEUE,
		"",	
		false,
		false,
		false,
		false,
		nil,
	)
}

func (mq *MQ) CloseConnection(){

	err := mq.channel.Close()
	if err != nil {
		util.FailOnError(err, "close channel:")
	}

	err = mq.connection.Close()
	if err != nil {
		util.FailOnError(err, "close connection:")
	}
}