package main

import (
	router "consumer/api/router"
	config "consumer/internal/config"
	mq "consumer/internal/queue"
	worker "consumer/pkg/worker"
	browser "consumer/pkg/browser"
	logger "consumer/pkg/log"

	"net/http"
	"time"
)

func main() {
	logger.Start()

	logger.Info("🧩 Configuring environment")
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load env %s", err)
	}

	logger.Info("📡 Connecting with RabbitMQ on port %s", cfg.RabbitMQ.Port)
	queue, err := mq.Create(&cfg.RabbitMQ)
	if err != nil {
		logger.Error("Failed to create connection with rabbitMQ %s", err)
	}
	
	defer queue.Close()

	logger.Info("📥 Initiating channel for consumption")
	jobs, err := queue.Consume()
	if err != nil {
		logger.Error("Failed to consume %s", err)
	}

	logger.Info("🧠 Creating and filling the browser pool")
	browsers, err := browser.CreatePool(3)
	if err != nil {
		logger.Error("Failed to start browsers %s", err)
	}

	defer browsers.Close()

	logger.Info("⚙️ Creating worker threads")
	go worker.Factory(jobs, browsers, 3, 3)

	server := http.Server{
		Addr: cfg.Server.Port,
		Handler: router.Config(queue, &cfg.Server),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("🚀 Starting server on port %s", cfg.Server.Port)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server %s", err)
	}

	logger.Info("🛑 Server exited")
}
