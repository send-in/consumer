package mq

import (
	"consumer/util"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Connection *amqp.Connection
var Channel *amqp.Channel

func CreateConnection(
	url string,
) (
	*amqp.Connection,
	*amqp.Channel,
) {
	// creating connection
	connection, err := amqp.Dial(url)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()
	
	// creating channel
	channel, err :=  connection.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	Connection, Channel = 
	connection, channel

	return connection, channel
}

func Send(

) (
	error,
) {
	if Connection == nil || Channel == nil{
		return errors.New("No Channel Found")
	}



}

func Consume(

) (
	error,
) {
	if Connection == nil || Channel == nil{
		return errors.New("No Channel Found")
	}
	

}