package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/jdnCreations/gcms/internal/database"
	"github.com/jdnCreations/gcms/internal/models"
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

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to create a user")
	decoder := json.NewDecoder(r.Body)
	params := models.UserInfo{} 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// validate user
	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
	}) 
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		respondWithError(w, 422, "Could not create user")
		return
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to retrieve all users")
	users, err := cfg.db.GetAllUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not retrieve users")
		return
	}

	respondWithJSON(w, http.StatusOK, users)	
}

func (cfg *apiConfig) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to delete user")
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

	err = cfg.db.DeleteUserById(context.Background(), uuid) 
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		respondWithError(w, http.StatusBadRequest, "Unable to delete user")
		return
	}

	respondWithJSON(w, http.StatusOK, "Deleted user")
}

func (cfg *apiConfig) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

	user, err := cfg.db.GetUserById(context.Background(), uuid)
	if err != nil {
		respondWithError(w, 404, "user not found")
		return
	}

	respondWithJSON(w, 200, user)
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := models.UpdateUserInfo{} 
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input data")
		return
	}

	cust, err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		Column1: params.FirstName,
		Column2: params.LastName,
		Column3: params.Email,
		ID: uuid,
	})
	if err != nil {
		respondWithError(w, 500, "Could not update user")
		return
	}

	respondWithJSON(w, http.StatusOK, cust)
}

func (cfg *apiConfig) handleGetAllGames(w http.ResponseWriter, r *http.Request) {
	games, err := cfg.db.GetAllGames(context.Background())
	if err != nil {
		log.Fatal(err)
		respondWithError(w, 500, "Could not retrieve all games")
	}

	respondWithJSON(w, 200, games)
}

func (cfg *apiConfig) handleCreateGame(w http.ResponseWriter, r *http.Request) {
  log.Println("Attempting to create game")
  decoder := json.NewDecoder(r.Body)
  gameInfo := models.GameInfo{}
	log.Println(gameInfo)
	err := decoder.Decode(&gameInfo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate game 
	validate := validator.New()
	err = validate.Struct(gameInfo)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	game, err := cfg.db.CreateGame(context.Background(), 
	database.CreateGameParams{
		Title: gameInfo.Title,
		Copies: gameInfo.Copies,
	})
	if err != nil {
		respondWithError(w, 422, "Could not create game, name must be unique")
		return
	}

  respondWithJSON(w, 200, game)
}

func (cfg *apiConfig) handleDeleteGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

	err = cfg.db.DeleteGameById(context.Background(), uuid)
	if err != nil {
		respondWithError(w, 404, "Game not found")
		return
	}

	respondWithJSON(w, 200, nil)
}

func (cfg *apiConfig) handleGetGameById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid game ID format")
		return
	}

	game, err := cfg.db.GetGameById(context.Background(), uuid)
	if err != nil {
		respondWithError(w, 404, "Game not found")
		return
	}

	respondWithJSON(w, 200, game)
}

func (cfg *apiConfig) handleUpdateGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid game ID format")
		return
	}

	gameToUpdate, err := cfg.db.GetGameById(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 404, "Game does not exist with that id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := models.UpdateGameInfo{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if params.Title != nil {
		gameToUpdate.Title = *params.Title
	}

	if params.Copies != nil {
		gameToUpdate.Copies = *params.Copies
	}
	
	game, err := cfg.db.UpdateGame(r.Context(), database.UpdateGameParams{
		Column1: gameToUpdate.Title,
		Column2: gameToUpdate.Copies,
		ID: gameToUpdate.ID,
	})
	if err != nil {
		respondWithError(w, 422, "Could not update game")
		return
	}

	respondWithJSON(w, 200, game)
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
		

		r := mux.NewRouter()	
		server := &http.Server{
			Addr: ":8080",
			Handler: r,
		}
		r.HandleFunc("/api/healthz", handleReadiness).Methods("GET")

		// users
		r.HandleFunc("/api/users", apiCfg.handleCreateUser).Methods("POST")
		r.HandleFunc("/api/users/{id}", apiCfg.handleGetUserById).Methods("GET")
		r.HandleFunc("/api/users", apiCfg.handleGetAllUsers).Methods("GET")
    r.HandleFunc("/api/users/{id}", apiCfg.handleUpdateUser).Methods("PUT")
		r.HandleFunc("/api/users/{id}", apiCfg.handleDeleteUser).Methods("DELETE")

		// games
		r.HandleFunc("/api/games", apiCfg.handleGetAllGames).Methods("GET")
		r.HandleFunc("/api/games", apiCfg.handleCreateGame).Methods("POST")
		r.HandleFunc("/api/games/{id}", apiCfg.handleGetGameById).Methods("GET")
		r.HandleFunc("/api/games/{id}", apiCfg.handleDeleteGame).Methods("DELETE")
		r.HandleFunc("/api/games/{id}", apiCfg.handleUpdateGame).Methods("PUT")

		// reservations
		r.HandleFunc("/api/reservations", apiCfg.handleActiveReservations).Methods("GET")
		server.ListenAndServe()
		
}