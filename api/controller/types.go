package controller

import mq "consumer/internal/queue"

type Jobs struct {
	queue *mq.MQ
}