package main

import (
	"log"
	"net/http"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/handler"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func main() {
	// create main service
	service := &handler.BackendService{}
	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	// create multiplexer to handle service and static file routes
	mux := http.NewServeMux()

	// mount the API
	mux.Handle("/", srv)

	// serve static files from "./static" under the "/static" route
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// start the hourly tasks in a separate goroutine
	go runHourlyTasks()

	// start the server
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func runHourlyTasks() {
	// initial run
	executeHourlyTasks()

	// create a ticker to run the tasks every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		executeHourlyTasks()
	}
}

func executeHourlyTasks() {
	est := utils.LoadEST()
	log.Println("Executing hourly tasks at:", time.Now().In(est))
}
