package main

import (
	"consumer/api/router"
	"consumer/lib"
	"consumer/mq"
	"consumer/util"
	
	"os"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/joho/godotenv"
)

type application struct{
	port string
	handler http.Handler
}

func main(){
	var err error

	// load env
	err = godotenv.Load(".env")
	util.FailOnError(err, "Failed to load env")
	
	// connect with rabbit mq
	log.Printf("Connecting with rabbit mq on port %s", mq.PORT,)
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", mq.USERNAME, mq.PASSWORD, mq.HOST, mq.PORT)
	queue, err := mq.CreateConnection(url)
	util.FailOnError(err, "Failed to create connection with rabbitMQ")
	defer queue.CloseConnection()

	requests, err := queue.Consume()
	util.FailOnError(err, "Failed to consume")
	lib.Workers(requests, 3)

	// initialize an application
	app := &application{
		port: ":" + os.Getenv("PORT"),
		handler: router.Handler(queue),
	}
	
	log.Printf("Starting server on port %s", app.port,)
	server := http.Server{
		Addr: app.port,
		Handler: app.handler,
		IdleTimeout: time.Minute,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 30*time.Second,
	}

	err = server.ListenAndServe()
	util.FailOnError(err, "Failed to start server")

	// exiting server
	log.Println("Server exiting")

	// message := mq.Message{
	// 	UserAgent: os.Getenv("TEST_AGENT"),
	// 	Token: os.Getenv("TEST_TOKEN"),
	// 	Message: "you never deserved ananya",
	// 	Reciever: "https://www.linkedin.com/in/visshon/",
	// }

	// browser := lib.CreateBrowser(&message)
	// defer browser.Cancel()

	// browser.Send(message)

	
}