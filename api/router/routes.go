package router

import (
	controller "consumer/api/controller"
	// middleware "consumer/api/middleware"

	mq "consumer/internal/queue"
	config "consumer/internal/config"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Config(mq *mq.MQ, cfg *config.ServerConfig) http.Handler {
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
	
	{
		jobs := controller.Create(mq)

		v1.GET(
			"/health", 
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "healthy",
				})
			},
		)
		
		v1.POST(
			"/jobs",
			// middleware.Authenticate(cfg.Passkey),
			jobs.POST,
		)
	}

	return router
}