package worker

import (
	browser "consumer/pkg/browser"
	
	"github.com/rabbitmq/amqp091-go"
)

type WorkerGroup struct {
	jobs <-chan amqp091.Delivery
	browsers *browser.BrowserPool
	count int
}