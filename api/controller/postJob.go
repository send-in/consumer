package controller

import (
	"consumer/mq"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (jobs *JobController) PostJob(context *gin.Context){
	queue := jobs.queue
	if queue == nil{
		context.JSON(
			http.StatusInternalServerError, 
			gin.H{"error": "No queue found"},
		)
        return
	}

	// auth check
	var message mq.Message

	err := context.ShouldBindJSON(&message); 
	if err != nil {
        context.JSON(
			http.StatusBadRequest, 
			gin.H{"error": err.Error()},
		)
        return
    }

	log.Println(message)
	queue.Publish(
		message,
	)
	
	context.JSON(
		http.StatusOK,
		gin.H{"data": "Golang RabbitMQ API POST"},
	)
}