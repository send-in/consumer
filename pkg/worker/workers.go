package worker

import (
	mq "consumer/internal/queue"

	browser "consumer/pkg/browser"
	logger "consumer/pkg/log"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

func Factory(jobs <-chan amqp.Delivery, browsers *browser.BrowserPool, maxThreads int, maxRetry int64) {

	var group errgroup.Group
	group.SetLimit(maxThreads)

	for request := range jobs {
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

		group.Go(func() error{
			jobid := request.MessageId
			Worker(jobid, request, browsers)
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		logger.Error("One or more workers failed: %v", err)
	}
}

func Worker(id string, request amqp.Delivery, pool *browser.BrowserPool) {
	var message mq.Message
	json.Unmarshal(request.Body, &message)
	
	logger.Info(
		"worker picked up the job %s", 
		id,
	)

	browser := pool.Acquire()
	status, err := browser.SendTemp()

	if err != nil {
		logger.Error(
			"[Job %s] Failed to send message, sending to dead queue: %v ", 
			id, err,
		)
		request.Nack(false, false)
	} else if status {
		logger.Success(
			"[Job %s] Message sent successfully", 
			id,
		)
		request.Ack(false)
	} else {
		logger.Info(
			"[Job %s] Message not confirmed, requeueing to main queue", 
			id,
		)
		request.Nack(false, true) 
	}

	pool.Release(browser)
}