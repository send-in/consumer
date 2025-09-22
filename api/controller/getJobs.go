package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (jobs *JobController) GetJobs(context *gin.Context) {
	context.JSON(
		http.StatusOK,
		gin.H{"data": "Golang RabbitMQ API GET"},
	)
}