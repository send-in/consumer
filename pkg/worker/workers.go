package worker

import (
	// mq "consumer/internal/queue"

	browser "consumer/pkg/browser"
	logger "consumer/pkg/log"

	// "encoding/json"
	// "time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"github.com/google/uuid"
)

func Factory(jobs <-chan amqp.Delivery, browsers *browser.BrowserPool, max int) {

	var group errgroup.Group
	group.SetLimit(max)

	for request := range jobs {
		group.Go(func() error{
			jobid := uuid.New().String()
			Worker(jobid, request, browsers)
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		logger.Error("One or more workers failed: %v", err)
	}
}

func Worker(id string, request amqp.Delivery, pool *browser.BrowserPool) {
	// var message mq.Message
	// json.Unmarshal(request.Body, &message)
	
	logger.Info(
		"worker picked up the job %s", 
		id,
	)

	browser := pool.Acquire()
	status, err := browser.SendTemp()

	if err != nil {
		logger.Error(
			"[Worker %s] Failed to send message, requeueing: %v ", 
			id, err,
		)
		request.Nack(false, true) // requeue
	} else if status {
		logger.Success(
			"[Worker %s] Message sent successfully", 
			id,
		)
		request.Ack(false)
	} else {
		logger.Info(
			"[Worker %s] Message not confirmed, requeueing", 
			id,
		)
		request.Nack(false, true) // requeue
	}

	pool.Release(browser)
}