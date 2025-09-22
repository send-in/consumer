package lib

import (
	"consumer/mq"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func Workers(requests <-chan amqp.Delivery, workerCount int) {
	for i := 1; i <= workerCount; i++ {
		go func(id int) {
			for d := range requests {
				switch d.Type {
					case mq.SENDMESSAGE:
						var msg mq.Message
						if err := json.Unmarshal(d.Body, &msg); err != nil {
							fmt.Printf("[Worker %d] Failed to unmarshal message: %v\n", id, err)
							d.Nack(false, false)
							continue
						}

						fmt.Printf("%s[Worker %d] 🚀 STARTED: UserAgent=%s, Message=%s%s\n",
							Yellow, id, msg.UserAgent, msg.Message, Reset,
						)

						time.Sleep(10 * time.Second)

						fmt.Printf("%s[Worker %d] ✅ DONE: UserAgent=%s, Message=%s%s\n",
							Green, id, msg.UserAgent, msg.Message, Reset,
						)
				}
			}
		}(i)
	}
}