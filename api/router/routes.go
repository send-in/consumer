package router

import (
	controller "consumer/api/controller"
	mq "consumer/internal/queue"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler(mq *mq.MQ) http.Handler {
	router := gin.Default()

	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST"},
			AllowHeaders: []string{"*"},
			AllowCredentials: true,
		},
	))

	v1 := router.Group("/api/v1")
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
	
	{
		jobs := controller.Create(mq)
		v1.GET("/jobs", jobs.GET)
		v1.POST("/jobs", jobs.POST)
	}

	return router
}