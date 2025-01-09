package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jdnCreations/gcms/internal/database"
	"github.com/jdnCreations/gcms/internal/models"
	"github.com/jdnCreations/gcms/internal/utils"
)

type apiConfig struct {
  db *database.Queries
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
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
	type User struct {
		FirstName string;
		LastName string;
		Email string;
	}
	log.Println("Attempting to create a user")
	decoder := json.NewDecoder(r.Body)
	params := User{} 
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
		log.Printf("Error: %s", err.Error())
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

	err = cfg.db.DeleteUserById(r.Context(), uuid) 
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

	type User struct {
		FirstName string
		LastName string
		Email string
	}

	decoder := json.NewDecoder(r.Body)
	params := User{} 
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
	games, err := cfg.db.GetAllGames(r.Context())
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

	game, err := cfg.db.CreateGame(r.Context(), 
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

	err = cfg.db.DeleteGameById(r.Context(), uuid)
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

func (cfg *apiConfig) handleCreateGenre(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  genreInfo := models.GenreInfo{}
  err := decoder.Decode(&genreInfo)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, err.Error())
		return
  }

  genre, err := cfg.db.CreateGenre(r.Context(), genreInfo.Name)
  if err != nil {
    respondWithError(w, 422, "Genre may already exist")
    return
  }

  respondWithJSON(w, 200, genre)
}

func (cfg *apiConfig) handleGetAllGenres(w http.ResponseWriter, r *http.Request) {
  genres, err := cfg.db.GetAllGenres(r.Context())
  if err != nil {
    respondWithError(w, 422, "Could not retrieve genres")
    return
  }

  respondWithJSON(w, 200, genres)
}

func (cfg *apiConfig) handleGetGenreById(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]
  uuid, err := uuid.Parse(id)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }

  genre, err := cfg.db.GetGenreById(r.Context(), uuid)
  if err != nil {
    respondWithError(w, 404, "No genre with that id")
    return
  }

  respondWithJSON(w, 200, genre)
}

func (cfg *apiConfig) handleDeleteGenre(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]
  uuid, err := uuid.Parse(id)
  if err != nil {
    respondWithError(w, 400, "invalid ID format")
    return
  }

  err = cfg.db.DeleteGenreById(r.Context(), uuid)
  if err != nil {
    respondWithError(w, 404, "No genre with that id")
    return
  }

  respondWithJSON(w, 200, "deleted")
}

func (cfg *apiConfig) handleUpdateGenre(w http.ResponseWriter, r *http.Request) {
  uuid, err := utils.GetIdFromRequest("id", r)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
  }

  genreToUpdate, err := cfg.db.GetGenreById(r.Context(), uuid)
  if err != nil {
    respondWithError(w, 404, "No genre with that id")
    return
  }

  type p struct {
    Name *string
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

  if params.Name != nil {
    genreToUpdate.Name = *params.Name
  }
  
  genre, err := cfg.db.UpdateGenre(r.Context(), database.UpdateGenreParams{
    Column1: genreToUpdate.Name,
    ID: genreToUpdate.ID,
  })
  if err != nil {
    respondWithError(w, 422, "Unable to update genre")
    return
  }

  respondWithJSON(w, 200, genre)
}

func (cfg *apiConfig) handleAddGenreToGame(w http.ResponseWriter, r *http.Request) {
  gameID, err := utils.GetIdFromRequest("gameID", r)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format - gameid")
    return
  }

  type p struct {
    GenreID string `json:"genre_id"`
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid request payload")
    return
  }

  genreID, err := uuid.Parse(params.GenreID)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format - genreID")
    return
  }

  gg, err := cfg.db.AddGenreToGame(r.Context(), database.AddGenreToGameParams{
    GameID: gameID,
    GenreID: genreID,
  })
  if err != nil {
    respondWithError(w, 422, "Could not add genre to game")
    return
  }

  respondWithJSON(w, 200, gg)
}

func (cfg *apiConfig) handleActiveReservations(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.db.GetAllActiveReservations(r.Context())
  if err != nil {
    respondWithError(w, 422, "Could not retrieve active reservations")
    return
  }

  respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleGetAllReservations(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.db.GetAllReservations(r.Context())
  if err != nil {
    respondWithError(w, 422, "Could not retrieve reservations")
    return
  }

  respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleCreateReservation(w http.ResponseWriter, r *http.Request) {
  type p struct {
    StartTime string
    EndTime string
    UserID string
    GameID string
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid request payload")
    return
  }

	log.Printf(params.StartTime)
	log.Printf(params.GameID)

	// convert to appropriate types for db ?
	startTime, err := utils.ConvertToPGTimestamp(params.StartTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid start time")
		return
	}

	endTime, err := utils.ConvertToPGTimestamp(params.EndTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid end time")
		return
	}
	
	userId, err := utils.ConvertToPGUUID(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid userid")
		return 
	}
	
	gameId, err := utils.ConvertToPGUUID(params.GameID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid gameid")
		return 
	}

  res, err := cfg.db.CreateReservation(r.Context(), database.CreateReservationParams{
    StartTime: startTime,
		EndTime: endTime,
		UserID: userId,
		GameID: gameId,
  })
	if err != nil {
		log.Printf("Error: %s", err.Error())
		respondWithError(w, 422, "Could not create reservation")
		return
	}

	respondWithJSON(w, 201, res)
}

func (cfg *apiConfig) handleGetReservationsForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	uuid, err := utils.ConvertToPGUUID(userID) 
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }

	res, err := cfg.db.GetReservationsForUser(r.Context(), uuid)
	if err != nil {
    respondWithError(w, 422, "Could not retrieve reservations")
    return
	}
	
  respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleCheckGameAvailable(w http.ResponseWriter, r *http.Request) {
	
}

func (cfg *apiConfig) handleDeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resID:= vars["id"]

	uuid, err := uuid.Parse(resID)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }

	err = cfg.db.DeleteReservation(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 422, "Could not delete reservation")
		return
	}

	respondWithJSON(w, 200, "")
}

func main() {
    err := godotenv.Load(".env")
    if err != nil {
      log.Fatalf("Could not load .env: %v", err)
    }
		port := os.Getenv("PORT")
		dbURL := os.Getenv("DATABASE_URL")
		if port == "" {
			log.Fatal("PORT environment variable is not set")
		}

		dbpool, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			log.Fatal(err)
		}
		defer dbpool.Close()

		// db, err := pgx.Connect(context.Background(), dbURL)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close(context.Background())

		dbQueries := database.New(dbpool)
		apiCfg := apiConfig{}
		apiCfg.db = dbQueries
		

		r := mux.NewRouter()	
		// wrap with cors and json content type
		corsRouter := enableCORS(jsonContentTypeMiddleware(r))
		server := &http.Server{
			Addr: ":8080",
			Handler: corsRouter,
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

    // genres
    r.HandleFunc("/api/genres", apiCfg.handleCreateGenre).Methods("POST")
    r.HandleFunc("/api/genres", apiCfg.handleGetAllGenres).Methods("GET")
    r.HandleFunc("/api/genres/{id}", apiCfg.handleGetGenreById).Methods("GET")
    r.HandleFunc("/api/genres/{id}", apiCfg.handleDeleteGenre).Methods("DELETE")
    r.HandleFunc("/api/genres/{id}", apiCfg.handleUpdateGenre).Methods("PUT")

    // add genre to games
    r.HandleFunc("/api/game_genres/{gameID}", apiCfg.handleAddGenreToGame).Methods("POST")

		// reservations
    r.HandleFunc("/api/reservations", apiCfg.handleCreateReservation).Methods("POST")
		r.HandleFunc("/api/reservations", apiCfg.handleGetAllReservations).Methods("GET")
		r.HandleFunc("/api/reservations/active", apiCfg.handleActiveReservations).Methods("GET")
		r.HandleFunc("/api/reservations/check/{gameID}", apiCfg.handleCheckGameAvailable).Methods("GET")
		r.HandleFunc("/api/reservations/{id}", apiCfg.handleDeleteReservation).Methods("DELETE")
		r.HandleFunc("/api/reservations/{userID}", apiCfg.handleGetReservationsForUser).Methods("GET")


		server.ListenAndServe()
		
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}