package main

import (
	router "consumer/internal/api/router"
	config "consumer/internal/config"
	worker "consumer/internal/worker"
	logger "consumer/pkg/log"
	mq "consumer/pkg/queue"

	"net/http"
	"time"
)

func main() {
	logger.Start()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load env %s", err)
	}

	logger.Success("Connecting with rabbit mq on port %s", cfg.RabbitMQ.Port)
	queue, err := mq.Create(&cfg.RabbitMQ)
	if err != nil {
		logger.Error("Failed to create connection with rabbitMQ %s", err)
	}

	defer queue.Close()

	requests, err := queue.Consume()
	if err != nil {
		logger.Error("Failed to consume %s", err)
	}
	
	logger.Info("Starting server on port %s", cfg.Server.Port)
	
	// worker.Workers(requests, 3, cfg.RabbitMQ.Type)
	
	server := http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router.Handler(queue),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server %s", err)
	}

	logger.Info("Server exited")
}
