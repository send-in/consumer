package controller

import "consumer/mq"

type JobController struct {
	queue *mq.MQ
}

func NewJobController(q *mq.MQ) *JobController {
	return &JobController{
		queue: q,
	}
}