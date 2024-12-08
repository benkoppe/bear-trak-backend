package main

import (
	"log"
	"net/http"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	handler "github.com/benkoppe/bear-trak-backend/handler"
)

func main() {
	service := &handler.BackendService{}
	srv, err := backend.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":3000", srv); err != nil {
		log.Fatal(err)
	}
}
