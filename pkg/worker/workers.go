package worker

import (
	mq "consumer/internal/queue"
	"context"

	browser "consumer/pkg/browser"
	logger "consumer/pkg/log"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

func Factory(
	jobs <-chan amqp.Delivery, 
	browsers *browser.BrowserPool, 
	maxThreads int, 
	maxRetry int64,
	context context.Context,
) {
	var group errgroup.Group
	group.SetLimit(maxThreads)

	for {
		select{
			case <-context.Done():
				closeWorkers(&group)
				return
			case request, ok := <-jobs:
				if !ok {
					closeWorkers(&group)
					return
				}

				death, exists := request.Headers["x-death"].([]any)

				if exists && len(death) > 0 {
					count, ok := death[0].(amqp.Table)["count"].(int64)
					if count >= maxRetry && ok {
						logger.Error("Max retries reached for message: %s", request.MessageId)
						request.Ack(false)

						// TODO: post request with message details
						continue
					}
				}

				message := request
				group.Go(func() error{
					Worker(message, browsers)
					return nil
				})
		}
	}
}

func Worker(
	request amqp.Delivery, 
	pool *browser.BrowserPool,
) {
	id := request.MessageId
	logger.Info("worker picked up the job %s", id)

	var message mq.Message
	err := json.Unmarshal(request.Body, &message)
	if err != nil {
		logger.Error("malformed object body: %s", err)
		request.Nack(false, false)
		return
	}
	
	browser := pool.Acquire()
	if browser == nil {
		logger.Error("Nil browser recieved for job: %s", id)
		request.Nack(false, true)
		return
	}

	status, err := browser.SendTemp()
	if err != nil {
		logger.Error("[Job %s] Failed to send message, sending to dead queue: %v ", id, err)
		request.Nack(false, false)
	} else if status {
		logger.Success("[Job %s] Message sent successfully", id)
		request.Ack(false)
	} else {
		logger.Info(
			"[Job %s] Message not confirmed, requeueing to main queue", id)
		request.Nack(false, true) 
	}

	pool.Release(browser)
}