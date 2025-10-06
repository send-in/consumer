package controller

import (
	logger "consumer/pkg/log"
	mq "consumer/pkg/queue"
	
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(q *mq.MQ) *Jobs {
	return &Jobs{
		queue: q,
	}
}

func (jobs *Jobs) GET(context *gin.Context) {
	context.JSON(
		http.StatusOK, 
		gin.H{"data": jobs.queue},
	)
}

func (jobs *Jobs) POST(context *gin.Context){
	var message mq.Message
	queue := jobs.queue

	if err := context.ShouldBindJSON(&message); err != nil {
		logger.Warning("Invalid request body: %v", err)
		context.JSON(
			http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := queue.Publish(message); err != nil {
		logger.Error("Failed to create job: %v", err)
		context.JSON(
			http.StatusInternalServerError, 
			gin.H{"error": "Failed to create job"},
		)
		return
	}
	
	context.JSON(
		http.StatusOK,
		gin.H{"data": message},
	)
}