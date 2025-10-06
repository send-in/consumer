package main

import (
	router "consumer/internal/api/router"
	config "consumer/internal/config"
	worker "consumer/internal/worker"
	browser "consumer/pkg/browser"
	logger "consumer/pkg/log"
	mq "consumer/pkg/queue"

	"net/http"
	"time"
)

func main() {
	logger.Start()

	logger.Info("[1] Configuring environment")
	cfg, err := config.Load()

	if err != nil {
		logger.Error("Failed to load env %s", err)
	}

	logger.Info("[2] Connecting with rabbit mq on port %s", cfg.RabbitMQ.Port)
	queue, err := mq.Create(&cfg.RabbitMQ)
	
	if err != nil {
		logger.Error("Failed to create connection with rabbitMQ %s", err)
	}
	
	defer queue.Close()

	logger.Info("[3] Initiating channel for consumtion")
	jobs, err := queue.Consume()

	if err != nil {
		logger.Error("Failed to consume %s", err)
	}

	logger.Info("[4] Creating and filling the browser pool")
	browsers, err := browser.CreatePool(1)

	if err != nil {
		logger.Error("Failed to start browsers %s", err)
	}

	logger.Info("[5] Creating worker threads")
	worker.Workers(jobs, browsers, 1)

	server := http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router.Handler(queue),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Starting server on port %s", cfg.Server.Port)
	err = server.ListenAndServe()

	if err != nil {
		logger.Error("Failed to start server %s", err)
	}

	logger.Info("[!] Server exited")
}
