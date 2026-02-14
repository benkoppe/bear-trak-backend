package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/bmizerany/pq"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms"
	"github.com/benkoppe/bear-trak-backend/go-server/schools"
)

//go:embed static/*
var embeddedStaticFiles embed.FS

func main() {
	// get school code from environment
	schoolCodeStr := os.Getenv("SCHOOL_CODE")
	if schoolCodeStr == "" {
		schoolCodeStr = string(schools.Cornell)
	}
	schoolCode := schools.SchoolCode(schoolCodeStr)
	config, err := schools.GetConfig(schoolCode)
	if err != nil {
		log.Fatalf("Error getting config for school %s: %v", schoolCode, err)
	}

	// connect to db
	pool, err := connectToDBPool(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer pool.Close()

	dbQueries := db.New(pool)

	// create main service
	handler, err := schools.NewHandler(schoolCode, dbQueries, config)
	if err != nil {
		log.Fatal(err)
	}
	srv, err := api.NewServer(handler)
	if err != nil {
		log.Fatal(err)
	}

	// allowed origins for domain cross-request
	allowedOrigins := []string{
		"http://localhost:3000",
		"https://trak.2ben.dev",
	}
	corsMiddleware := corsMiddleware(allowedOrigins)

	// create multiplexer to handle service and static file routes
	mux := http.NewServeMux()

	// mount the API
	mux.Handle("/", corsMiddleware(srv))

	// create a sub filesystem rooted as "static"
	staticFS, err := fs.Sub(embeddedStaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	// serve static files from the embedded filesystem on /static/
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", corsMiddleware(http.StripPrefix("/static", fileServer)))

	// serve internal export route if token is set
	internalExportToken := os.Getenv("INTERNAL_EXPORT_TOKEN")
	if internalExportToken == "" {
		log.Println("INTERNAL_EXPORT_TOKEN is not set; internal capacities export route is disabled")
	} else {
		mux.Handle(
			"/internal/v1/gyms/capacities/export",
			gyms.NewCapacitiesExportHandler(dbQueries, internalExportToken),
		)
	}

	// start the timed tasks in a separate goroutine
	go runTimedTasks(dbQueries, handler, *config)

	// start the server
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func runTimedTasks(queries *db.Queries, handler api.Handler, config schools.Config) {
	// initial run
	executeTimedTasks(queries, handler, config)

	// create a ticker to run the tasks
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		executeTimedTasks(queries, handler, config)
	}
}

func executeTimedTasks(queries *db.Queries, handler api.Handler, config schools.Config) {
	// est := timeutils.LoadEST()
	// log.Println("Executing timed tasks at:", time.Now().In(est))

	// create a context with a timeout for the operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if config.EnabledGymCapacities {
		err := gyms.LogCapacities(ctx, handler, queries)
		if err != nil {
			log.Printf("Error logging gym capacities: %v", err)
		}
	}
}

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowed[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
