package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/KMalkowski/rss-feed/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("POSTGRES"))
	if err != nil {
		panic("can't establish database connection")
	}

	api := &apiConfig{
		DB: database.New(db),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", HealthzHandler)
	mux.HandleFunc("GET /v1/err", ErrHandler)
	mux.HandleFunc("POST /v1/users", api.HandleCreateUser)
	mux.HandleFunc("GET /v1/users", api.authMiddleware(api.GetUserByApiKeyHandler))
	mux.HandleFunc("POST /v1/feeds", api.authMiddleware(api.CreateFeedHandler))
	mux.HandleFunc("GET /v1/feeds", api.GetFeedsHandler)
	mux.HandleFunc("POST /v1/feed_follows", api.authMiddleware(api.CreateFeedFollowHandler))

	server := http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	log.Fatal(server.ListenAndServe())
	server.ListenAndServe()
	log.Println("serving on the port " + os.Getenv("PORT"))
}
