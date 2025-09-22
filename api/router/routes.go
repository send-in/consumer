package router

import (
	"consumer/api/controller"
	
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler() http.Handler {
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

	{
		v1.GET("/jobs", controller.GetJobs)
		v1.POST("/jobs", controller.PostJobs)
	}

	return httpRouter
}