package worker

import (
	logger "consumer/pkg/log"
	mq "consumer/pkg/queue"

	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func Workers(requests <-chan amqp091.Delivery, count int, Type string) {
	for i := 1; i <= count; i++ {
		go func(id int) {
			for d := range requests {
				switch d.Type {
					case Type:{
						var msg mq.Message

						if err := json.Unmarshal(d.Body, &msg); err != nil {
							fmt.Printf(
								"[Worker %d] Failed to unmarshal message: %v\n", 
								id, err,
							)
							d.Nack(false, false)
							continue
						}

						logger.Info(
							"Worker %d started processing: UserAgent=%s, Message=%s",
							id, msg.UserAgent, msg.Message,
						)

						time.Sleep(10 * time.Second)

						logger.Success(
							"Worker %d completed processing: UserAgent=%s, Message=%s",
							id, msg.UserAgent, msg.Message,
						)
					}
				}
			}
		}(i)	
	}
}