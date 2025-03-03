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
	"github.com/benkoppe/bear-trak-backend/go-server/gyms"
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

	dbQueries := db.New(pool)

	// create main service
	service := &handler.BackendService{
		DB: dbQueries,
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

	// start the timed tasks in a separate goroutine
	go runTimedTasks(dbQueries)

	// start the server
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func runTimedTasks(queries *db.Queries) {
	// initial run
	executeHourlyTasks(queries)

	// create a ticker to run the tasks
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		executeHourlyTasks(queries)
	}
}

func executeHourlyTasks(queries *db.Queries) {
	est := utils.LoadEST()
	log.Println("Executing timed tasks at:", time.Now().In(est))

	// create a context with a timeout for the operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := gyms.LogCapacities(ctx, handler.GymCapacitiesUrl, queries)
	if err != nil {
		log.Printf("Error logging gym capacities: $v", err)
	}
}
