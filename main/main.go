package main

import (
	"log"
	"net/http"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/handler"
)

func main() {
	// create main service
	service := &handler.BackendService{}
	srv, err := backend.NewServer(service)
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

	// start the server
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
