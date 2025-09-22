package main

import (
	"consumer/api/router"
	"consumer/util"

	"log"
	"net/http"
	"os"
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
	
	// initialize an application
	app := &application{
		port: ":" + os.Getenv("PORT"),
		handler: router.Handler(),
	}
	
	server := http.Server{
		Addr: app.port,
		Handler: app.handler,
		IdleTimeout: time.Minute,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 30*time.Second,
	}

	log.Printf("Starting server on port %s", app.port,)
	
	// starting server
	err = server.ListenAndServe()
	util.FailOnError(err, "Failed to start server")


	// exiting server
	log.Println("Server exiting")
}