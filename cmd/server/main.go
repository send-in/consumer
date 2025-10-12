package main

import (
	router "consumer/api/router"
	config "consumer/internal/config"
	mq "consumer/internal/queue"
	browser "consumer/pkg/browser"
	logger "consumer/pkg/log"
	worker "consumer/pkg/worker"

	"context"
	"net/http"
	"time"

	"github.com/ory/graceful"
)

func main() {
	logger.Start()

	logger.Info("🧩 Configuring environment")
	cfg, err := config.Load()
	logger.Fatal(err, "Failed to load env")

	logger.Info("📡 Connecting with RabbitMQ on port %s", cfg.RabbitMQ.Port)
	queue, err := mq.Create(&cfg.RabbitMQ)
	logger.Fatal(err, "Failed to create connection with rabbitMQ")
	defer queue.Close()

	logger.Info("📥 Initiating channel for consumption")
	jobs, err := queue.Consume()
	logger.Fatal(err, "Failed to consume")

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("🧠 Creating and filling the browser pool")
	browsers, err := browser.CreatePool(
		3, // max instances
		context,
	)
	logger.Fatal(err, "Failed to start browsers")
	defer browsers.Close()

	logger.Info("⚙️ Creating worker threads")
	go worker.Factory(
		jobs, 
		browsers, 
		3, // max threads
		3, // max retry
		context,
	)

	server := graceful.WithDefaults(
		&http.Server{
			Addr:         cfg.Server.Port,
			Handler:      router.Config(queue, &cfg.Server),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	)

	logger.Info("🚀 Starting server on port %s", cfg.Server.Port)
	err = graceful.Graceful(
		server.ListenAndServe, 
		server.Shutdown,
	)
	logger.Fatal(err, "Failed to gracefully shutdown")

	logger.Info("🛑 Server exited")
}
