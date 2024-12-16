package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jdnCreations/gcms/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
		// run http server 
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
		}

		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("PORT environment variable is not set")
		}

		apiCfg := apiConfig{}

		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			log.Println("DATABASE_URL environment variable is not set")
			log.Println("Running without CRUD endpoints")
		} else {
			db, err := sql.Open("postgres", dbURL)
			if err != nil {
				log.Fatal(err)
			}
			dbQueries := database.New(db)
			apiCfg.DB = dbQueries
			log.Println("Connected to database!")
		}

		mux := http.NewServeMux()	
		server := &http.Server{
			Addr: ":8080",
			Handler: mux,
		}
		mux.HandleFunc("GET /api/healthz", handleReadiness)
		server.ListenAndServe()
		
}