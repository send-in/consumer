package worker

import (
	logger "consumer/pkg/log"
	browser "consumer/pkg/browser"
	mq "consumer/pkg/queue"

	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
)


func Workers(jobs <-chan amqp091.Delivery, browsers *browser.BrowserPool ,count int) {
	pool := &WorkerGroup{
		jobs,
		browsers,
		count,
	}

	for id := range count {
		go pool.Worker(id)
	}
}

func (workers *WorkerGroup) Worker(id int) {
	for request := range workers.jobs {
		var message mq.Message

		if err := json.Unmarshal(request.Body, &message); err != nil {
			logger.Error(
				"[Worker %d] Failed to unmarshal message: %v\n", 
				id, err,
			)
			request.Nack(false, false) // dont requeue
			continue
		}
		
		logger.Info("worker picked up the job %d, for %s", id, message.Message)
		workers.Do(
			id, 
			message, 
			request,
		)
	}
}

func (workers *WorkerGroup) Do(id int, message mq.Message, request amqp091.Delivery) {
	pool := workers.browsers
	browser := pool.Acquire()

	if browser == nil {
		logger.Error(
			"[Worker %d] No available browser, skipping message, requeueing", 
			id,
		)
		request.Nack(false, true) // requeue
		return
	}

	status, err := browser.Send(message)

	if err != nil {
		logger.Error(
			"[Worker %d] Failed to send message, requeueing: %v ", 
			id, err,
		)
		request.Nack(false, true) // requeue
	} else if status {
		logger.Success(
			"[Worker %d] Message sent successfully", 
			id,
		)
		request.Ack(false)
	} else {
		logger.Info(
			"[Worker %d] Message not confirmed, requeueing", 
			id,
		)
		request.Nack(false, true) // requeue
	}
}