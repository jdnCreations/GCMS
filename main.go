package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/jdnCreations/gcms/internal/database"
	"github.com/jdnCreations/gcms/internal/models"
	"github.com/jdnCreations/gcms/internal/utils"
	"github.com/joho/godotenv"
)

type apiConfig struct {
  db *database.Queries
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handleActiveReservations(w http.ResponseWriter, r *http.Request) {

}


func respondWithError(w http.ResponseWriter, code int, msg string) {
  type returnErr struct {
    Error string `json:"error"`
  }

  respBody := returnErr{
    Error: msg,
  }
  dat, err := json.Marshal(respBody)
  if err != nil {
    log.Printf("Error marshalling JSON %s", err)
    w.WriteHeader(500)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  dat, err := json.Marshal(payload)
  if err != nil {
    respondWithError(w, 500, "Error marshalling JSON")
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(dat)
}

func (cfg *apiConfig) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to create a customer")
	decoder := json.NewDecoder(r.Body)
	params := models.CustomerInfo{} 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	missingFields := []string{}
	if params.FirstName == "" {
		missingFields = append(missingFields, "first_name")
	}
	if params.LastName == "" {
		missingFields = append(missingFields, "last_name")
	}
	if params.Email == "" {
		missingFields = append(missingFields, "email")
	}

	if len(missingFields) > 0 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Missing required fields: %s", strings.Join(missingFields, ", ")))
		return
	}

	// validate email
	err = utils.ValidateEmail(params)
	if err != nil {
		respondWithError(w, 422, "Invalid email")
		return
	}

	user, err := cfg.db.CreateCustomer(r.Context(), database.CreateCustomerParams{
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
	}) 
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		respondWithError(w, 422, "Could not create customer")
		return
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to retrieve all customers")
	customers, err := cfg.db.GetAllCustomers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not retrieve customers")
		return
	}

	respondWithJSON(w, http.StatusOK, customers)	
}

func (cfg *apiConfig) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to delete customer")
  decoder := json.NewDecoder(r.Body)
	var params struct {
		ID string `json:"id"`
	}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if params.ID == "" {
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	uuid, err := uuid.Parse(params.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID format")
		return
	}

	err = cfg.db.DeleteCustomerById(context.Background(), uuid) 
	if err != nil {
		log.Printf("Failed to delete customer: %v", err)
		respondWithError(w, http.StatusBadRequest, "Unable to delete customer")
		return
	}

	respondWithJSON(w, http.StatusOK, "Deleted customer")
}



func main() {
		// run http server 
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
		}
		
		port := os.Getenv("PORT")
		dbURL := os.Getenv("DATABASE_URL")
		if port == "" {
			log.Fatal("PORT environment variable is not set")
		}

		
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		
		dbQueries := database.New(db)
		apiCfg := apiConfig{}
		apiCfg.db = dbQueries
		

		mux := http.NewServeMux()	
		server := &http.Server{
			Addr: ":8080",
			Handler: mux,
		}
		mux.HandleFunc("GET /api/healthz", handleReadiness)
		mux.HandleFunc("GET /api/reservations", apiCfg.handleActiveReservations)
		mux.HandleFunc("POST /api/customers", apiCfg.handleCreateCustomer)
		mux.HandleFunc("GET /api/customers", apiCfg.handleGetAllCustomers)
    // mux.HandleFunc("PUT /api/customers", apiCfg.handleUpdateCustomer)
		mux.HandleFunc("DELETE /api/customers", apiCfg.handleDeleteCustomer)
		server.ListenAndServe()
		
}