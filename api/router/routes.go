package router

import (
	"consumer/api/controller"
	"consumer/mq"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler(mq *mq.MQ) http.Handler {
	httpRouter := gin.Default()

	httpRouter.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST"},
			AllowHeaders: []string{"*"},
			AllowCredentials: true,
		},
	))

	v1 := httpRouter.Group("/api/v1")

	jobsController := controller.NewJobController(mq)

	{
		v1.GET("/jobs", jobsController.GetJobs)
		v1.POST("/jobs", jobsController.PostJob)
	}

	return httpRouter
}