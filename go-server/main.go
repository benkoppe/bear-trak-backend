package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	_ "github.com/bmizerany/pq"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/handler"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

//go:embed static/*
var embeddedStaticFiles embed.FS

func main() {
	// connect to db
	pool, err := connectToDbPool(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer pool.Close()

	// create main service
	service := &handler.BackendService{
		DB: db.New(pool),
	}
	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}

	// create multiplexer to handle service and static file routes
	mux := http.NewServeMux()

	// mount the API
	mux.Handle("/", srv)

	// create a sub filesystem rooted as "static"
	staticFS, err := fs.Sub(embeddedStaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	// serve static files from the embedded filesystem on /static/
	fileServer := http.FileServer(http.FS(staticFS))
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
