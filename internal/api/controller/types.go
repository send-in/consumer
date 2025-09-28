package controller

import mq "consumer/pkg/queue"

type Jobs struct {
	queue *mq.MQ
}